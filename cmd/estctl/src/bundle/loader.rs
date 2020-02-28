use super::*;
use crate::specs::canonical;
use std::collections::HashSet;
use std::io;
use thiserror;

#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("IO error: {0}")]
    Io(#[from] io::Error),
    #[error("failed to parse URL: {0}")]
    URLErr(#[from] url::ParseError),
    #[error("failed to parse YAML: {0}")]
    YAMLErr(#[from] serde_yaml::Error),
    #[error("failed to parse JSON: {0}")]
    JSONErr(#[from] serde_json::Error),
    #[error("failed to compile JSON-Schema: {0}")]
    BuildErr(#[from] schema::build::Error),
    #[error("bundle database error: {0}")]
    SQLiteErr(#[from] rusqlite::Error),
    #[error("schema index error: {0}")]
    IndexErr(#[from] schema::index::Error),
    #[error("converting from SQL: {0}")]
    SQLConversionErr(#[from] rusqlite::types::FromSqlError),
    #[error("{0}")]
    CanonicalError(#[from] specs::canonical::Error),
}
use Error::*;
use rusqlite::OptionalExtension;

pub trait FileSystem {
    fn open(&self, url: &url::Url) -> Result<Box<dyn io::Read>, io::Error>;
}

pub struct Loader {
    fs: Box<dyn FileSystem>,
    db: rusqlite::Connection,
}

type Schema = schema::Schema<specs::Annotation>;

impl Loader {

    pub fn define_schema(db: rusqlite::Connection) -> rusqlite::Result<rusqlite::Connection> {
        db.execute_batch("
        ")?;

        Ok(db)
    }

    pub fn new(db: rusqlite::Connection, fs: Box<dyn FileSystem>) -> Result<Loader, Error> {
        let db = Loader::define_schema(db)?;

        db.execute_batch("
            ATTACH ':memory:' as tmp;
            CREATE TABLE tmp.seen ( url TEXT NOT NULL PRIMARY KEY );

            BEGIN;
            PRAGMA defer_foreign_keys = true;
        ")?;

        Ok(Loader {
            fs,
            db,
        })
    }

    pub fn build_schema_catalog(&self) -> Result<Vec<Schema>, Error> {
        let mut schemas = self.db.prepare_cached("SELECT url, body FROM schema;")?;
        let schemas = schemas.
            query_and_then(rusqlite::NO_PARAMS, |row: &rusqlite::Row| -> Result<Schema, Error> {
                let url = url::Url::parse(row.get_raw(0).as_str()?)?;
                let raw = serde_json::from_str(row.get_raw(1).as_str()?)?;
                Ok(schema::build::build_schema(url, &raw)?)
            })?;

        let mut catalog = Vec::new();
        for schema in schemas {
            catalog.push(schema?);
        }
        Ok(catalog)
    }

    pub fn finish(self) -> Result<rusqlite::Connection, Error> {
        let catalog = self.build_schema_catalog()?;

        let mut ind = schema::index::Index::new();
        for s in &catalog {
            ind.add(s)?;
        }
        ind.verify_references()?;

        /*
        let tmp: Option<String> = self.db.query_row(
            "SELECT url FROM collection WHERE schema_url NOT IN (SELECT url FROM schema);",
            rusqlite::NO_PARAMS,
            |row| row.get(0),
        ).optional()?;

        if let Some(url) = tmp {
            panic!("got url {} ", url)
        };
        */

        self.db.execute_batch("
            COMMIT;
            DETACH DATABASE tmp;
        ")?;

        Ok(self.db)
    }

    fn already_seen(&mut self, url: &url::Url) -> Result<bool, Error> {
        let mut s = self.db.prepare_cached(
            "INSERT INTO tmp.seen (url) VALUES (?) ON CONFLICT DO NOTHING;")?;

        Ok(s.execute(&[url.as_str()])? == 0)
    }

    pub fn load_node(&mut self, node: url::Url) -> Result<(), Error> {
        if self.already_seen(&node)? {
            return Ok(())
        }
        println!("loading {}", &node);

        let br = io::BufReader::new(self.fs.open(&node)?);
        let spec: specs::Node = serde_yaml::from_reader(br)?;
        //let spec = spec.into_canonical(&node)?;

        for inc in &spec.include {
            let inc = canonical::join(&node, &inc)?;
            self.load_node(inc)?;
        }

        for c in &spec.collections {
            let name = canonical::join(&node, &c.name)?;

            let schema = canonical::join(&node, &c.schema)?;
            self.load_schema(schema.clone())?;

            let mut s = self.db.prepare_cached(
                "INSERT INTO collection (url, schema_url, key, partitions)\
                        VALUES (?, ?, ?, ?);")?;

            s.execute(&[
                name.as_str(),
                schema.as_str(),
                serde_json::to_string(&c.key)?.as_str(),
                serde_json::to_string(&c.partitions)?.as_str(),
                ])?;
        }

        Ok(())
    }

    fn load_schema(&mut self, mut url: url::Url) -> Result<(), Error> {
        url.set_fragment(None);

        if self.already_seen(&url)? {
            return Ok(())
        }
        println!("loading {}", url);

        let br = io::BufReader::new(self.fs.open(&url)?);
        let raw_schema: serde_json::Value = {
            if url.path().ends_with(".yaml") {
                serde_yaml::from_reader(br)?
            } else {
                serde_json::from_reader(br)?
            }
        };

        let compiled: Schema = schema::build::build_schema(url.clone(), &raw_schema)?;
        self.walk_schema(&url, &compiled)?;

        let mut s = self.db.prepare_cached(
            "INSERT INTO schema (url, body) VALUES (?, ?);")?;

        s.execute(&[url.as_str(), raw_schema.to_string().as_str()])?;
        Ok(())
    }

    fn walk_schema(&mut self, base: &url::Url, schema: &Schema) -> Result<(), Error> {
        use schema::Keyword;
        use schema::Application;

        Ok(for kw in &schema.kw {
            match kw {
                Keyword::Application(Application::Ref(ruri), _)  => {
                    self.load_schema(ruri.clone())?
                },
                Keyword::Application(_, child) => self.walk_schema(base, &child)?,
                // No-ops.
                Keyword::Anchor(_) | Keyword::RecursiveAnchor | Keyword::Validation(_) | Keyword::Annotation(_) => (),
            }
        })
    }

    /*
    fn process_root(&mut self, base: url::Url, spec: specs::Project) -> Result<(), Error> {
        for mut c in spec.collections {
            c.name = base.join(&c.name)?.to_string();
            c.schema = base.join(&c.schema)?.to_string();

            if !c.examples.is_empty() {
                c.examples = base.join(&c.examples)?.to_string();
            }
            if let Some(d) = &mut c.derivation {
                self.process_derivation(&base, d)
            }
        }
        Ok(())
    }

    fn process_derivation(&mut self, base: &url::Url, spec: &mut specs::Derivation) -> Result<(), Error> {
        use specs::Derivation;

        match &mut d {
            Derivation::Jq(d) =>
        }
    }

    fn process_path(&mut self, base: &url::Url, path: &mut String) -> Result<(), Error> {
        Ok(*path = base.join(path)?.to_string())
    }

    */
}
