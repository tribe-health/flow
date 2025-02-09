use super::references::*;

use lazy_static::lazy_static;
use regex::Regex;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};
use std::time::Duration;
use validator::Validate;

#[derive(Serialize, Deserialize, Clone, Debug, JsonSchema, Validate)]
#[schemars(example = "BucketAndPrefix::example")]
pub struct BucketAndPrefix {
    /// Bucket into which Flow will store data.
    #[validate(regex = "BUCKET_RE")]
    bucket: String,

    /// Optional prefix of keys written to the bucket.
    #[validate]
    #[serde(default)]
    prefix: Option<Prefix>,
}

impl BucketAndPrefix {
    fn bucket_and_prefix(&self) -> (&str, &str) {
        (self.bucket.as_str(), self.prefix.as_deref().unwrap_or(""))
    }

    pub fn example() -> Self {
        Self {
            bucket: "my-bucket".to_string(),
            prefix: None,
        }
    }
}

/// Details of an s3-compatible storage endpoint, such as Minio or R2.
#[derive(Serialize, Deserialize, Clone, Debug, JsonSchema, Validate)]
#[schemars(example = "CustomStore::example")]
pub struct CustomStore {
    /// Bucket into which Flow will store data.
    #[validate(regex = "BUCKET_RE")]
    pub bucket: String,
    /// endpoint is required when provider is "custom", and specifies the
    /// address of an s3-compatible storage provider.
    pub endpoint: StorageEndpoint,
    /// Optional prefix of keys written to the bucket.
    #[validate]
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub prefix: Option<Prefix>,
}

impl CustomStore {
    pub fn example() -> Self {
        Self {
            bucket: "my-bucket".to_string(),
            endpoint: StorageEndpoint::example(),
            prefix: None,
        }
    }

    fn bucket_and_prefix(&self) -> (&str, &str) {
        (self.bucket.as_str(), self.prefix.as_deref().unwrap_or(""))
    }
}

/// A Store into which Flow journal fragments may be written.
///
/// The persisted path of a journal fragment is determined by composing the
/// Store's bucket and prefix with the journal name and a content-addressed
/// fragment file name.
///
/// Eg, given a Store to S3 with bucket "my-bucket" and prefix "a/prefix",
/// along with a collection "example/events" having a logical partition "region",
/// then a complete persisted path might be:
///
///   s3://my-bucket/a/prefix/example/events/region=EU/utc_date=2021-10-25/utc_hour=13/000123-000456-789abcdef.gzip
///
#[derive(Serialize, Deserialize, Clone, Debug, JsonSchema)]
#[schemars(example = "Store::example")]
#[serde(tag = "provider", rename_all = "SCREAMING_SNAKE_CASE")]
pub enum Store {
    ///# Amazon Simple Storage Service.
    S3(BucketAndPrefix),
    ///# Google Cloud Storage.
    Gcs(BucketAndPrefix),
    ///# Azure object storage service.
    Azure(BucketAndPrefix),
    ///# An S3-compatible endpoint
    Custom(CustomStore),
}

impl Validate for Store {
    fn validate(&self) -> Result<(), validator::ValidationErrors> {
        match self {
            Self::S3(s) | Self::Gcs(s) | Self::Azure(s) => s.validate(),
            Self::Custom(s) => s.validate(),
        }
    }
}

impl Store {
    pub fn example() -> Self {
        Self::S3(BucketAndPrefix::example())
    }
    pub fn to_url(&self, catalog_name: &str) -> url::Url {
        let (scheme, (bucket, prefix)) = match self {
            Self::S3(cfg) => ("s3", cfg.bucket_and_prefix()),
            Self::Gcs(cfg) => ("gs", cfg.bucket_and_prefix()),
            Self::Azure(cfg) => ("azure", cfg.bucket_and_prefix()),
            // Custom storage endpoints are expected to be s3-compatible, and thus use the s3 scheme
            Self::Custom(cfg) => ("s3", cfg.bucket_and_prefix()),
        };
        let mut url = url::Url::parse(&format!("{}://{}/{}", scheme, bucket, prefix))
            .expect("parsing as URL should never fail");
        if let Store::Custom(cfg) = self {
            let tenant = catalog_name
                .split_once('/')
                .expect("invalid catalog_name passed to Store::to_url")
                .0;
            url.query_pairs_mut()
                .append_pair("profile", tenant)
                .append_pair("endpoint", &cfg.endpoint);
        }
        url
    }
}

/// Storage defines the backing cloud storage for journals.
#[derive(Serialize, Deserialize, Clone, Debug, JsonSchema, Validate)]
pub struct StorageDef {
    /// # Stores for journal fragments under this prefix.
    ///
    /// Multiple stores may be specified, and all stores are periodically scanned
    /// to index applicable journal fragments. New fragments are always persisted
    /// to the first store in the list.
    ///
    /// This can be helpful in performing bucket migrations: adding a new store
    /// to the front of the list causes ongoing data to be written to that location,
    /// while historical data continues to be read and served from the prior stores.
    ///
    /// When running `flowctl test`, stores are ignored and a local temporary
    /// directory is used instead.
    #[validate]
    pub stores: Vec<Store>,
}

impl StorageDef {
    pub fn example() -> Self {
        Self {
            stores: vec![Store::example()],
        }
    }
}

/// A CompressionCodec may be applied to compress journal fragments before
/// they're persisted to cloud stoage. The compression applied to a journal
/// fragment is included in its filename, such as ".gz" for GZIP. A
/// collection's compression may be changed at any time, and will affect
/// newly-written journal fragments.
#[derive(Deserialize, Debug, Serialize, JsonSchema, Clone)]
#[serde(deny_unknown_fields, rename_all = "SCREAMING_SNAKE_CASE")]
#[schemars(example = "CompressionCodec::example")]
pub enum CompressionCodec {
    None,
    Gzip,
    Zstandard,
    Snappy,
    GzipOffloadDecompression,
}

impl CompressionCodec {
    pub fn example() -> Self {
        CompressionCodec::GzipOffloadDecompression
    }
}

/// A FragmentTemplate configures how journal fragment files are
/// produced as part of a collection.
// path_postfix_template and refresh_interval are deliberately not
// exposed here. We're fixing these values in place for now.
#[derive(Serialize, Deserialize, Debug, Default, JsonSchema, Validate, Clone)]
#[serde(deny_unknown_fields, rename_all = "camelCase")]
#[schemars(example = "FragmentTemplate::example")]
pub struct FragmentTemplate {
    /// # Desired content length of each fragment, in megabytes before compression.
    /// When a collection journal fragment reaches this threshold, it will be
    /// closed off and pushed to cloud storage.
    /// If not set, a default of 512MB is used.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    #[validate(range(min = 32, max = 4096))]
    pub length: Option<u32>,
    /// # Codec used to compress Journal Fragments.
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub compression_codec: Option<CompressionCodec>,
    /// # Duration for which historical fragments of a collection should be kept.
    /// If not set, then fragments are retained indefinitely.
    #[serde(
        default,
        with = "humantime_serde",
        skip_serializing_if = "Option::is_none"
    )]
    #[schemars(schema_with = "super::duration_schema")]
    pub retention: Option<std::time::Duration>,
    /// # Maximum flush delay before in-progress fragments are closed and persisted
    /// into cloud storage. Intervals are converted into uniform time segments:
    /// 24h will "roll" all fragments at midnight UTC every day, 1h at the top of
    /// every hour, 15m a :00, :15, :30, :45 past the hour, and so on.
    /// If not set, then fragments are not flushed on time-based intervals.
    #[serde(
        default,
        with = "humantime_serde",
        skip_serializing_if = "Option::is_none"
    )]
    #[schemars(schema_with = "super::duration_schema")]
    pub flush_interval: Option<std::time::Duration>,
}

impl FragmentTemplate {
    pub fn example() -> Self {
        Self {
            compression_codec: Some(CompressionCodec::Zstandard),
            flush_interval: Some(Duration::from_secs(3600)),
            ..Default::default()
        }
    }
    pub fn is_empty(&self) -> bool {
        let FragmentTemplate {
            length: o1,
            compression_codec: o2,
            retention: o3,
            flush_interval: o4,
        } = self;

        o1.is_none() && o2.is_none() && o3.is_none() && o4.is_none()
    }
}

/// A JournalTemplate configures the journals which make up the
/// physical partitions of a collection.
#[derive(Serialize, Deserialize, Debug, Default, JsonSchema, Clone)]
#[serde(deny_unknown_fields, rename_all = "camelCase")]
#[schemars(example = "JournalTemplate::example")]
pub struct JournalTemplate {
    /// # Fragment configuration of collection journals.
    pub fragments: FragmentTemplate,
}

impl JournalTemplate {
    pub fn example() -> Self {
        Self {
            fragments: FragmentTemplate::example(),
        }
    }
    pub fn is_empty(&self) -> bool {
        let JournalTemplate { fragments } = self;
        fragments.is_empty()
    }
}

lazy_static! {
    // BUCKET_RE matches a cloud provider bucket. Simplified from (look-around removed):
    // https://stackoverflow.com/questions/50480924/regex-for-s3-bucket-name
    static ref BUCKET_RE: Regex =
        Regex::new(r#"(^(([a-z0-9]|[a-z0-9][a-z0-9\-]*[a-z0-9])\.)*([a-z0-9]|[a-z0-9][a-z0-9\-]*[a-z0-9])$)"#).unwrap();
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_regexes() {
        for (case, expect) in [
            ("foo.bar.baz", true),
            ("foo-bar-baz", true),
            ("foo/bar/baz", false),
            ("Foo.Bar.Baz", false),
        ] {
            assert!(BUCKET_RE.is_match(case) == expect);
        }
    }

    // The representation of Store was changed from a struct to an enum, so this test is ensuring
    // that existing Stores will deserialize properly with the new representation.
    #[test]
    fn old_store_json_still_deserializes_into_new_enum() {
        let actual: Store =
            serde_json::from_str(r#"{"provider":"GCS","prefix":"flow/","bucket":"test-bucket"}"#)
                .expect("failed to deserialize");
        let Store::Gcs(b_and_p) = actual else {
            panic!("expected a gcs store, got: {:?}", actual);
        };
        assert_eq!("test-bucket", &b_and_p.bucket);
        assert_eq!(Some("flow/"), b_and_p.prefix.as_deref());
    }

    #[test]
    fn custom_storage_endpoint() {
        let actual: Store = serde_json::from_str(
            r#"{"provider":"CUSTOM","prefix":"test/","bucket":"test-bucket", "endpoint": "http://canary.test:1234"}"#,
        ).expect("failed to deserialize");
        let Store::Custom(cfg) = &actual else {
            panic!("expected a custom store, got: {:?}", actual);
        };
        assert_eq!("http://canary.test:1234", cfg.endpoint.as_str());
        assert_eq!("test-bucket", &cfg.bucket);
        assert_eq!(Some("test/"), cfg.prefix.as_deref());

        actual.validate().expect("failed validation");

        let actual_url = actual.to_url("testTenant/foo").to_string();
        assert_eq!(
            "s3://test-bucket/test/?profile=testTenant&endpoint=http%3A%2F%2Fcanary.test%3A1234",
            &actual_url
        );
    }
}
