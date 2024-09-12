use proto::meta::v1::meta_client::MetaClient;
use proto::payload::v1::meta;
use prost_types::Any;
use opentelemetry::global;
use opentelemetry::propagation::Injector;
use tonic::metadata::MetadataMap;
use tonic::Request;

// MetadataInjector構造体を定義して、Injectorトレイトを実装する
struct MetadataInjector<'a>(&'a mut MetadataMap);

impl<'a> Injector for MetadataInjector<'a> {
    fn set(&mut self, key: &str, value: String) {
        // keyを所有権を持つ文字列に変換してから、挿入する
        let key_owned = key.to_owned(); // Stringに変換
        let parsed_key = key_owned.parse::<tonic::metadata::MetadataKey<_>>().unwrap();
        self.0.insert(parsed_key, value.parse().unwrap());
    }
}

// ジェネリック型を用いたトレースコンテキストの挿入関数
fn inject_trace_context<T>(request: &mut Request<T>) {
    let metadata = request.metadata_mut();  // mutは不要
    let mut injector = MetadataInjector(metadata);  // injectorをミュータブルに

    // Contextの取得
    let current_context = opentelemetry::Context::current();

    // グローバルプロパゲータを使ってトレースコンテキストを挿入
    global::get_text_map_propagator(|propagator| {
        propagator.inject_context(&current_context, &mut injector);  // injectorをミュータブルで渡す
    });
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // gRPC サーバーに接続するクライアントを作成
    let mut client = MetaClient::connect("http://[::1]:8081").await?;

    // 1. set key: aaa, value: 1
    let any_value = Any {
        type_url: "".to_string(),  // 型URLを設定（必要なら適切なURLに）
        value: b"1".to_vec(),      // Vec<u8> を prost_types::Any に変換
    };
    let mut request = tonic::Request::new(meta::KeyValue {
        key: Some(meta::Key{key: "aaa".to_string()}),
        value: Some(meta::Value { value: Some(any_value) }),
    });
    // トレースコンテキストをリクエストに挿入
    inject_trace_context(&mut request);
    client.set(request).await?;
    println!("Set key: aaa, value: 1");

    // 2. set key: bbb, value: 2
    let any_value = Any {
        type_url: "".to_string(),  // 型URLを設定
        value: b"2".to_vec(),
    };
    let mut request = tonic::Request::new(meta::KeyValue {
        key: Some(meta::Key{key: "bbb".to_string()}),
        value: Some(meta::Value { value: Some(any_value) }),
    });
    inject_trace_context(&mut request); // トレースコンテキストを挿入
    client.set(request).await?;
    println!("Set key: bbb, value: 2");

    // 3. get key: bbb
    let mut request = tonic::Request::new(meta::Key {
        key: "bbb".to_string(),
    });
    inject_trace_context(&mut request); // トレースコンテキストを挿入
    let response = client.get(request).await?;
    println!("Get key: bbb, RESPONSE={:?}", response.into_inner());

    // 4. set key: bbb, value: 3
    let any_value = Any {
        type_url: "".to_string(),
        value: b"3".to_vec(),
    };
    let mut request = tonic::Request::new(meta::KeyValue {
        key: Some(meta::Key{key: "bbb".to_string()}),
        value: Some(meta::Value { value: Some(any_value) }),
    });
    inject_trace_context(&mut request); // トレースコンテキストを挿入
    client.set(request).await?;
    println!("Set key: bbb, value: 3");

    // 5. get key: bbb
    let mut request = tonic::Request::new(meta::Key {
        key: "bbb".to_string(),
    });
    inject_trace_context(&mut request); // トレースコンテキストを挿入
    let response = client.get(request).await?;
    println!("Get key: bbb, RESPONSE={:?}", response.into_inner());

    // 6. delete key: aaa
    let mut request = tonic::Request::new(meta::Key {
        key: "aaa".to_string(),
    });
    inject_trace_context(&mut request); // トレースコンテキストを挿入
    client.delete(request).await?;
    println!("Deleted key: aaa");

    // 7. get key: aaa (after deletion)
    let mut request = tonic::Request::new(meta::Key {
        key: "aaa".to_string(),
    });
    inject_trace_context(&mut request); // トレースコンテキストを挿入
    let response = client.get(request).await;
    match response {
        Ok(res) => println!("Get key: aaa, RESPONSE={:?}", res.into_inner()),
        Err(e) => println!("Get key: aaa failed with error: {:?}", e),
    }

    Ok(())
}

// #[tokio::main]
// async fn main() -> Result<(), Box<dyn std::error::Error>> {
//     // gRPC サーバーに接続するクライアントを作成
//     let mut client = MetaClient::connect("http://[::1]:8081").await?;

//     // 1. set key: aaa, value: 1
//     let any_value = Any {
//         type_url: "".to_string(),  // 型URLを設定（必要なら適切なURLに）
//         value: b"1".to_vec(),      // Vec<u8> を prost_types::Any に変換
//     };
//     let request = tonic::Request::new(meta::KeyValue {
//         key: Some(meta::Key{key: "aaa".to_string()}),
//         value: Some(meta::Value { value: Some(any_value) }),
//     });
//     client.set(request).await?;
//     println!("Set key: aaa, value: 1");

//     // 2. set key: bbb, value: 2
//     let any_value = Any {
//         type_url: "".to_string(),  // 型URLを設定
//         value: b"2".to_vec(),
//     };
//     let request = tonic::Request::new(meta::KeyValue {
//         key: Some(meta::Key{key: "bbb".to_string()}),
//         value: Some(meta::Value { value: Some(any_value) }),
//     });
//     client.set(request).await?;
//     println!("Set key: bbb, value: 2");

//     // 3. get key: bbb
//     let request = tonic::Request::new(meta::Key {
//         key: "bbb".to_string(),
//     });
//     let response = client.get(request).await?;
//     println!("Get key: bbb, RESPONSE={:?}", response.into_inner());

//     // 4. set key: bbb, value: 3
//     let any_value = Any {
//         type_url: "".to_string(),
//         value: b"3".to_vec(),
//     };
//     let request = tonic::Request::new(meta::KeyValue {
//         key: Some(meta::Key{key: "bbb".to_string()}),
//         value: Some(meta::Value { value: Some(any_value) }),
//     });
//     client.set(request).await?;
//     println!("Set key: bbb, value: 3");

//     // 5. get key: bbb
//     let request = tonic::Request::new(meta::Key {
//         key: "bbb".to_string(),
//     });
//     let response = client.get(request).await?;
//     println!("Get key: bbb, RESPONSE={:?}", response.into_inner());

//     // 6. delete key: aaa
//     let request = tonic::Request::new(meta::Key {
//         key: "aaa".to_string(),
//     });
//     client.delete(request).await?;
//     println!("Deleted key: aaa");

//     // 7. get key: aaa (after deletion)
//     let request = tonic::Request::new(meta::Key {
//         key: "aaa".to_string(),
//     });
//     let response = client.get(request).await;
//     match response {
//         Ok(res) => println!("Get key: aaa, RESPONSE={:?}", res.into_inner()),
//         Err(e) => println!("Get key: aaa failed with error: {:?}", e),
//     }

//     Ok(())
// }