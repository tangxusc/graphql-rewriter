# graphql-rewriter

从http协议中接受graphql查询,并重写其中的查询

关键流程:
```mermaid
graph LR
    query --graphql--> graphql-parser
    subgraph rewriter-pipeline
    graphql-parser --json schema--> wasm-A 
    graphql-parser --json schema--> wasm-B
    graphql-parser --json schema--> wasm-X
    wasm-A --json merge path--> graphqlFormater
    wasm-B --json merge path--> graphqlFormater
    wasm-X --json merge path--> graphqlFormater
    end
    graphqlFormater --format--new graphql--> dgraph[(graphql db)]
```

主要技术栈:
- gin
- graphql
- gqlparser
- json-patch
- wasmedge
- wasmdege-bindgen

主要功能:

- 替换部分查询