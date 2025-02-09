use super::Format;
use std::collections::BTreeMap;

// Indirect sub-locations within `sources` into external resources which
// are referenced through relative imports.
pub fn indirect_large_files(sources: &mut tables::Sources, threshold: usize) {
    let tables::Sources {
        captures,
        collections,
        fetches: _,
        imports: _,
        materializations,
        resources,
        storage_mappings: _,
        tests,
        errors: _,
    } = sources;

    for capture in captures.iter_mut() {
        indirect_capture(capture, resources, threshold);
    }
    for collection in collections.iter_mut() {
        indirect_collection(collection, resources, threshold);
    }
    for materialization in materializations.iter_mut() {
        indirect_materialization(materialization, resources, threshold);
    }
    for test in tests.iter_mut() {
        indirect_test(test, resources, threshold);
    }
}

// Extend Resources with Resource instances for each catalog specification
// URL which is referenced by any and all imports, captures, collections,
// materializations, and tests.
pub fn rebuild_catalog_resources(sources: &mut tables::Sources) {
    let tables::Sources {
        captures,
        collections,
        fetches: _,
        imports,
        materializations,
        resources,
        storage_mappings: _,
        tests,
        errors: _,
    } = sources;

    let mut catalogs: BTreeMap<url::Url, models::Catalog> = BTreeMap::new();

    let strip_scope = |scope: &url::Url| {
        let mut scope = scope.clone();
        scope.set_fragment(None);
        scope
    };

    for tables::Import { scope, to_resource } in imports.iter() {
        if !scope.fragment().unwrap().starts_with("/import") {
            continue; // Skip implicit imports.
        }
        let scope = strip_scope(scope);
        let import = match scope.make_relative(&to_resource) {
            Some(rel) => rel,
            None => to_resource.to_string(),
        };

        let entry = catalogs.entry(scope).or_default();
        entry.import.push(models::RelativeUrl::new(import));
    }

    for tables::Capture {
        scope,
        capture,
        spec,
    } in captures.iter()
    {
        let entry = catalogs.entry(strip_scope(scope)).or_default();
        entry.captures.insert(capture.clone(), spec.clone());
    }

    for tables::Collection {
        scope,
        collection,
        spec,
    } in collections.iter()
    {
        let entry = catalogs.entry(strip_scope(scope)).or_default();
        entry.collections.insert(collection.clone(), spec.clone());
    }

    for tables::Materialization {
        scope,
        materialization,
        spec,
    } in materializations.iter()
    {
        let entry = catalogs.entry(strip_scope(scope)).or_default();
        entry
            .materializations
            .insert(materialization.clone(), spec.clone());
    }

    for tables::Test { scope, test, spec } in tests.iter() {
        let entry = catalogs.entry(strip_scope(scope)).or_default();
        entry.tests.insert(test.clone(), spec.clone());
    }

    for (resource, mut catalog) in catalogs {
        catalog.import.sort();
        catalog.import.dedup();

        let content_dom: models::RawValue =
            serde_json::value::to_raw_value(&catalog).unwrap().into();
        let content_raw = Format::from_scope(&resource).serialize(&content_dom);

        tables::Resource {
            resource,
            content_dom,
            content: content_raw,
            content_type: proto_flow::flow::ContentType::Catalog,
        }
        .upsert_if_changed(resources)
    }
}

fn indirect_capture(
    capture: &mut tables::Capture,
    resources: &mut tables::Resources,
    threshold: usize,
) {
    let tables::Capture {
        scope,
        capture,
        spec: models::CaptureDef {
            endpoint, bindings, ..
        },
    } = capture;
    let base = base_name(capture);

    match endpoint {
        models::CaptureEndpoint::Connector(models::ConnectorConfig { config, .. }) => {
            indirect_dom(
                scope,
                config,
                format!("{base}.config"),
                resources,
                threshold,
            );
        }
    }

    for (index, models::CaptureBinding { resource, .. }) in bindings.iter_mut().enumerate() {
        indirect_dom(
            scope,
            resource,
            format!("{base}.resource.{index}.config"),
            resources,
            threshold,
        )
    }
}

fn indirect_collection(
    collection: &mut tables::Collection,
    resources: &mut tables::Resources,
    threshold: usize,
) {
    let tables::Collection {
        scope,
        collection,
        spec:
            models::CollectionDef {
                schema,
                write_schema,
                read_schema,
                key: _,
                projections: _,
                journals: _,
                derive,
                derivation: _,
            },
    } = collection;
    let base = base_name(collection);

    if let Some(schema) = schema {
        indirect_schema(
            scope,
            schema,
            format!("{base}.schema"),
            resources,
            threshold,
        );
    }
    if let Some(write_schema) = write_schema {
        indirect_schema(
            scope,
            write_schema,
            format!("{base}.write.schema"),
            resources,
            threshold,
        )
    }
    if let Some(read_schema) = read_schema {
        indirect_schema(
            scope,
            read_schema,
            format!("{base}.read.schema"),
            resources,
            threshold,
        );
    }
    if let Some(derivation) = derive {
        indirect_derivation(scope, derivation, base, resources, threshold);
    }
}

fn indirect_derivation(
    scope: &url::Url,
    derivation: &mut models::Derivation,
    base: &str,
    resources: &mut tables::Resources,
    threshold: usize,
) {
    let models::Derivation {
        using,
        transforms,
        shuffle_key_types: _,
        shards: _,
    } = derivation;
    let mut is_sql = false;

    match using {
        models::DeriveUsing::Connector(models::ConnectorConfig { config, .. }) => {
            indirect_dom(
                scope,
                config,
                format!("{base}.config"),
                resources,
                threshold,
            );
        }
        models::DeriveUsing::Sqlite(models::DeriveUsingSqlite { migrations }) => {
            is_sql = true;

            for (index, migration) in migrations.iter_mut().enumerate() {
                indirect_raw(
                    scope,
                    migration,
                    format!("{base}.migration.{index}.sql"),
                    resources,
                    threshold,
                );
            }
        }
        models::DeriveUsing::Typescript(models::DeriveUsingTypescript { module }) => {
            indirect_raw(scope, module, format!("{base}.ts"), resources, threshold);
        }
    }

    for models::TransformDef {
        name,
        lambda,
        shuffle,
        ..
    } in transforms
    {
        if is_sql {
            indirect_raw(
                scope,
                lambda,
                format!("{base}.lambda.{name}.sql"),
                resources,
                threshold,
            );
            if let models::Shuffle::Lambda(lambda) = shuffle {
                indirect_raw(
                    scope,
                    lambda,
                    format!("{base}.lambda.{name}.shuffle.sql"),
                    resources,
                    threshold,
                );
            }
        } else {
            indirect_dom(
                scope,
                lambda,
                format!("{base}.lambda.{name}"),
                resources,
                threshold,
            );
            if let models::Shuffle::Lambda(lambda) = shuffle {
                indirect_dom(
                    scope,
                    lambda,
                    format!("{base}.lambda.{name}.shuffle"),
                    resources,
                    threshold,
                );
            }
        }
    }
}

fn indirect_materialization(
    materialization: &mut tables::Materialization,
    resources: &mut tables::Resources,
    threshold: usize,
) {
    let tables::Materialization {
        scope,
        materialization,
        spec: models::MaterializationDef {
            endpoint, bindings, ..
        },
    } = materialization;
    let base = base_name(materialization);

    match endpoint {
        models::MaterializationEndpoint::Connector(models::ConnectorConfig { config, .. }) => {
            indirect_dom(
                scope,
                config,
                format!("{base}.config"),
                resources,
                threshold,
            )
        }
        _ => {}
    }

    for (index, models::MaterializationBinding { resource, .. }) in bindings.iter_mut().enumerate()
    {
        indirect_dom(
            scope,
            resource,
            format!("{base}.resource.{index}.config"),
            resources,
            threshold,
        )
    }
}

fn indirect_test(test: &mut tables::Test, resources: &mut tables::Resources, threshold: usize) {
    let tables::Test { scope, test, spec } = test;
    let base = base_name(test);

    for (index, step) in spec.iter_mut().enumerate() {
        let documents = match step {
            models::TestStep::Ingest(models::TestStepIngest { documents, .. })
            | models::TestStep::Verify(models::TestStepVerify { documents, .. }) => documents,
        };
        indirect_dom(
            scope,
            documents,
            format!("{base}.step.{index}"),
            resources,
            threshold,
        );
    }
}

fn indirect_schema(
    scope: &url::Url,
    content_dom: &mut models::RawValue,
    filename: String,
    resources: &mut tables::Resources,
    threshold: usize,
) {
    let schema = content_dom.to_value();

    // Attempt to clean up the schema by removing a superfluous $id.
    match schema {
        serde_json::Value::Object(mut m) => {
            if m.contains_key("definitions") || m.contains_key("$defs") {
                // We can't touch $id, as it provides the canonical base against which
                // $ref is resolved to definitions.
            } else if let Some(true) = m
                .get("$id")
                .and_then(serde_json::Value::as_str)
                .map(|s| s.starts_with("file://"))
            {
                m.remove("$id");
                *content_dom = models::RawValue::from_value(&serde_json::Value::Object(m))
            }
        }
        _ => (),
    };

    indirect_dom(scope, content_dom, filename, resources, threshold)
}

fn indirect_dom(
    scope: &url::Url,
    content_dom: &mut models::RawValue,
    filename: String,
    resources: &mut tables::Resources,
    threshold: usize,
) {
    if content_dom.get().len() <= threshold {
        // Leave small DOMs in-place.
        // This includes content_dom's which are already indirect.
        return;
    }

    let fmt = Format::from_scope(scope);
    let filename = format!("{filename}.{}", fmt.extension());

    tables::Resource {
        resource: scope.join(&filename).unwrap(),
        content_type: proto_flow::flow::ContentType::Config,
        content: fmt.serialize(content_dom),
        content_dom: content_dom.clone(),
    }
    .upsert_if_changed(resources);

    *content_dom =
        models::RawValue::from_string(serde_json::to_string(&filename).unwrap()).unwrap();
}

fn indirect_raw(
    scope: &url::Url,
    content_dom: &mut models::RawValue,
    filename: String,
    resources: &mut tables::Resources,
    threshold: usize,
) {
    if content_dom.get().len() <= threshold {
        // Leave small raw strings in-place.
        // This includes content_dom's which are already indirect.
        return;
    }

    let content_str =
        serde_json::from_str::<String>(content_dom.get()).expect("value must be a JSON string");

    tables::Resource {
        resource: scope.join(&filename).unwrap(),
        content_type: proto_flow::flow::ContentType::Config,
        content: content_str.into(),
        content_dom: std::mem::take(content_dom),
    }
    .upsert_if_changed(resources);

    *content_dom =
        models::RawValue::from_string(serde_json::to_string(&filename).unwrap()).unwrap();
}

fn base_name(name: &impl AsRef<str>) -> &str {
    let name = name.as_ref();

    match name.rsplit_once("/") {
        Some((_, base)) => base,
        None => name,
    }
}
