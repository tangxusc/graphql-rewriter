# graphql-rewriter

从http协议中接受graphql查询,并重写其中的查询

关键流程:
```mermaid
graph LR
    query --graphql--> graphql-parser
    subgraph rewriter-pipeline
    graphql-parser --> plugin-A 
    graphql-parser --> plugin-B
    graphql-parser --> plugin-X
    plugin-A --> graphqlFormater
    plugin-B --> graphqlFormater
    plugin-X --> graphqlFormater
    end
    graphqlFormater --format--new graphql--> dgraph[(graphql db)]
```

主要技术栈:
- gin
- graphql
- gqlparser
- ~~json-patch~~
- ~~wasmedge~~
- ~~wasmdege-bindgen~~
- golang plugin

主要功能:

- 替换部分查询