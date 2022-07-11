"use strict";(self.webpackChunksite=self.webpackChunksite||[]).push([[8623],{3905:function(e,t,n){n.d(t,{Zo:function(){return d},kt:function(){return c}});var a=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function l(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?l(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):l(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function o(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},l=Object.keys(e);for(a=0;a<l.length;a++)n=l[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(e);for(a=0;a<l.length;a++)n=l[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var p=a.createContext({}),s=function(e){var t=a.useContext(p),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},d=function(e){var t=s(e.components);return a.createElement(p.Provider,{value:t},e.children)},m={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},u=a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,l=e.originalType,p=e.parentName,d=o(e,["components","mdxType","originalType","parentName"]),u=s(n),c=r,k=u["".concat(p,".").concat(c)]||u[c]||m[c]||l;return n?a.createElement(k,i(i({ref:t},d),{},{components:n})):a.createElement(k,i({ref:t},d))}));function c(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var l=n.length,i=new Array(l);i[0]=u;var o={};for(var p in t)hasOwnProperty.call(t,p)&&(o[p]=t[p]);o.originalType=e,o.mdxType="string"==typeof e?e:r,i[1]=o;for(var s=2;s<l;s++)i[s]=n[s];return a.createElement.apply(null,i)}return a.createElement.apply(null,n)}u.displayName="MDXCreateElement"},2498:function(e,t,n){n.r(t),n.d(t,{assets:function(){return d},contentTitle:function(){return p},default:function(){return c},frontMatter:function(){return o},metadata:function(){return s},toc:function(){return m}});var a=n(7462),r=n(3366),l=(n(7294),n(3905)),i=["components"],o={sidebar_position:2},p="Amazon S3",s={unversionedId:"reference/Connectors/capture-connectors/amazon-s3",id:"reference/Connectors/capture-connectors/amazon-s3",title:"Amazon S3",description:"This connector captures data from an Amazon S3 bucket.",source:"@site/docs/reference/Connectors/capture-connectors/amazon-s3.md",sourceDirName:"reference/Connectors/capture-connectors",slug:"/reference/Connectors/capture-connectors/amazon-s3",permalink:"/reference/Connectors/capture-connectors/amazon-s3",draft:!1,editUrl:"https://github.com/estuary/flow/edit/master/site/docs/reference/Connectors/capture-connectors/amazon-s3.md",tags:[],version:"current",sidebarPosition:2,frontMatter:{sidebar_position:2},sidebar:"tutorialSidebar",previous:{title:"Amazon Kinesis",permalink:"/reference/Connectors/capture-connectors/amazon-kinesis"},next:{title:"Amplitude",permalink:"/reference/Connectors/capture-connectors/amplitude"}},d={},m=[{value:"Prerequisites",id:"prerequisites",level:2},{value:"Configuration",id:"configuration",level:2},{value:"Properties",id:"properties",level:3},{value:"Endpoint",id:"endpoint",level:4},{value:"Bindings",id:"bindings",level:4},{value:"Sample",id:"sample",level:3},{value:"Advanced: Parsing cloud storage data",id:"advanced-parsing-cloud-storage-data",level:3}],u={toc:m};function c(e){var t=e.components,n=(0,r.Z)(e,i);return(0,l.kt)("wrapper",(0,a.Z)({},u,n,{components:t,mdxType:"MDXLayout"}),(0,l.kt)("h1",{id:"amazon-s3"},"Amazon S3"),(0,l.kt)("p",null,"This connector captures data from an Amazon S3 bucket."),(0,l.kt)("p",null,(0,l.kt)("a",{parentName:"p",href:"https://ghcr.io/estuary/source-s3:dev"},(0,l.kt)("inlineCode",{parentName:"a"},"ghcr.io/estuary/source-s3:dev"))," provides the latest connector image. You can also follow the link in your browser to see past image versions."),(0,l.kt)("h2",{id:"prerequisites"},"Prerequisites"),(0,l.kt)("p",null,"To use this connector, either your S3 bucket must be public,\nor you must have access via a root or ",(0,l.kt)("a",{parentName:"p",href:"https://docs.aws.amazon.com/IAM/latest/UserGuide/id_users.html"},"IAM user"),"."),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},"For public buckets, verify that the ",(0,l.kt)("a",{parentName:"li",href:"https://docs.aws.amazon.com/AmazonS3/latest/userguide/access-control-overview.html#access-control-resources-manage-permissions-basics"},"access policy")," allows anonymous reads."),(0,l.kt)("li",{parentName:"ul"},"For buckets accessed by a user account, you'll need the AWS ",(0,l.kt)("strong",{parentName:"li"},"access key")," and ",(0,l.kt)("strong",{parentName:"li"},"secret access key")," for the user.\nSee the ",(0,l.kt)("a",{parentName:"li",href:"https://aws.amazon.com/blogs/security/wheres-my-secret-access-key/"},"AWS blog")," for help finding these credentials.")),(0,l.kt)("h2",{id:"configuration"},"Configuration"),(0,l.kt)("p",null,"You configure connectors either in the Flow web app, or by directly editing the catalog spec YAML.\nSee ",(0,l.kt)("a",{parentName:"p",href:"/concepts/connectors#using-connectors"},"connectors")," to learn more about using connectors. The values and YAML sample below provide configuration details specific to the S3 source connector."),(0,l.kt)("div",{className:"admonition admonition-tip alert alert--success"},(0,l.kt)("div",{parentName:"div",className:"admonition-heading"},(0,l.kt)("h5",{parentName:"div"},(0,l.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,l.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"12",height:"16",viewBox:"0 0 12 16"},(0,l.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M6.5 0C3.48 0 1 2.19 1 5c0 .92.55 2.25 1 3 1.34 2.25 1.78 2.78 2 4v1h5v-1c.22-1.22.66-1.75 2-4 .45-.75 1-2.08 1-3 0-2.81-2.48-5-5.5-5zm3.64 7.48c-.25.44-.47.8-.67 1.11-.86 1.41-1.25 2.06-1.45 3.23-.02.05-.02.11-.02.17H5c0-.06 0-.13-.02-.17-.2-1.17-.59-1.83-1.45-3.23-.2-.31-.42-.67-.67-1.11C2.44 6.78 2 5.65 2 5c0-2.2 2.02-4 4.5-4 1.22 0 2.36.42 3.22 1.19C10.55 2.94 11 3.94 11 5c0 .66-.44 1.78-.86 2.48zM4 14h5c-.23 1.14-1.3 2-2.5 2s-2.27-.86-2.5-2z"}))),"tip")),(0,l.kt)("div",{parentName:"div",className:"admonition-content"},(0,l.kt)("p",{parentName:"div"},"You might organize your S3 bucket using ",(0,l.kt)("a",{parentName:"p",href:"https://docs.aws.amazon.com/AmazonS3/latest/userguide/using-prefixes.html"},"prefixes")," to emulate a directory structure.\nThis connector can use prefixes in two ways: first, to perform the ",(0,l.kt)("a",{parentName:"p",href:"/concepts/connectors#flowctl-discover"},(0,l.kt)("strong",{parentName:"a"},"discovery"))," phase of setup, and later, when the capture is running."),(0,l.kt)("ul",{parentName:"div"},(0,l.kt)("li",{parentName:"ul"},"You can specify a prefix in the endpoint configuration to limit the overall scope of data discovery."),(0,l.kt)("li",{parentName:"ul"},"You're required to specify prefixes on a per-binding basis. This allows you to map each prefix to a distinct Flow collection,\nand informs how the capture will behave in production.")),(0,l.kt)("p",{parentName:"div"},"To capture the entire bucket, omit ",(0,l.kt)("inlineCode",{parentName:"p"},"prefix")," in the endpoint configuration and set ",(0,l.kt)("inlineCode",{parentName:"p"},"stream")," to the name of the bucket."))),(0,l.kt)("h3",{id:"properties"},"Properties"),(0,l.kt)("h4",{id:"endpoint"},"Endpoint"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Property"),(0,l.kt)("th",{parentName:"tr",align:null},"Title"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"),(0,l.kt)("th",{parentName:"tr",align:null},"Type"),(0,l.kt)("th",{parentName:"tr",align:null},"Required/Default"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/ascendingKeys")),(0,l.kt)("td",{parentName:"tr",align:null},"Ascending Keys"),(0,l.kt)("td",{parentName:"tr",align:null},"Improve sync speeds by listing files from the end of the last sync, rather than listing the entire bucket prefix. This requires that you write objects in ascending lexicographic order, such as an RFC-3339 timestamp, so that key ordering matches modification time ordering."),(0,l.kt)("td",{parentName:"tr",align:null},"boolean"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"false"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/awsAccessKeyId")),(0,l.kt)("td",{parentName:"tr",align:null},"AWS Access Key ID"),(0,l.kt)("td",{parentName:"tr",align:null},"Part of the AWS credentials that will be used to connect to S3. Required unless the bucket is public and allows anonymous listings and reads."),(0,l.kt)("td",{parentName:"tr",align:null},"string"),(0,l.kt)("td",{parentName:"tr",align:null})),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/awsSecretAccessKey")),(0,l.kt)("td",{parentName:"tr",align:null},"AWS Secret Access Key"),(0,l.kt)("td",{parentName:"tr",align:null},"Part of the AWS credentials that will be used to connect to S3. Required unless the bucket is public and allows anonymous listings and reads."),(0,l.kt)("td",{parentName:"tr",align:null},"string"),(0,l.kt)("td",{parentName:"tr",align:null})),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("strong",{parentName:"td"},(0,l.kt)("inlineCode",{parentName:"strong"},"/bucket"))),(0,l.kt)("td",{parentName:"tr",align:null},"Bucket"),(0,l.kt)("td",{parentName:"tr",align:null},"Name of the S3 bucket"),(0,l.kt)("td",{parentName:"tr",align:null},"string"),(0,l.kt)("td",{parentName:"tr",align:null},"Required")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/endpoint")),(0,l.kt)("td",{parentName:"tr",align:null},"AWS Endpoint"),(0,l.kt)("td",{parentName:"tr",align:null},"The AWS endpoint URI to connect to. Use if you","'","re capturing from a S3-compatible API that isn","'","t provided by AWS"),(0,l.kt)("td",{parentName:"tr",align:null},"string"),(0,l.kt)("td",{parentName:"tr",align:null})),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/matchKeys")),(0,l.kt)("td",{parentName:"tr",align:null},"Match Keys"),(0,l.kt)("td",{parentName:"tr",align:null},"Filter applied to all object keys under the prefix. If provided, only objects whose absolute path matches this regex will be read. For example, you can use ",'"',".","*","\\",".json",'"'," to only capture json files."),(0,l.kt)("td",{parentName:"tr",align:null},"string"),(0,l.kt)("td",{parentName:"tr",align:null})),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser")),(0,l.kt)("td",{parentName:"tr",align:null},"Parser Configuration"),(0,l.kt)("td",{parentName:"tr",align:null},"Configures how files are parsed"),(0,l.kt)("td",{parentName:"tr",align:null},"object"),(0,l.kt)("td",{parentName:"tr",align:null})),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/compression")),(0,l.kt)("td",{parentName:"tr",align:null},"Compression"),(0,l.kt)("td",{parentName:"tr",align:null},"Determines how to decompress the contents. The default, ","'","Auto","'",", will try to determine the compression automatically."),(0,l.kt)("td",{parentName:"tr",align:null},"null, string"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"null"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Determines how to parse the contents. The default, ","'","Auto","'",", will try to determine the format automatically based on the file extension or MIME type, if available."),(0,l.kt)("td",{parentName:"tr",align:null},"object"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},'{"auto":{}}'))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format/auto")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"object"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"{}"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format/avro")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"object"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"{}"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format/csv")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"object"),(0,l.kt)("td",{parentName:"tr",align:null})),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format/csv/delimiter")),(0,l.kt)("td",{parentName:"tr",align:null},"Delimiter"),(0,l.kt)("td",{parentName:"tr",align:null},"The delimiter that separates values within each row. Only single-byte delimiters are supported."),(0,l.kt)("td",{parentName:"tr",align:null},"null, string"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"null"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format/csv/encoding")),(0,l.kt)("td",{parentName:"tr",align:null},"Encoding"),(0,l.kt)("td",{parentName:"tr",align:null},"The character encoding of the source file. If unspecified, then the parser will make a best-effort guess based on peeking at a small portion of the beginning of the file. If known, it is best to specify. Encodings are specified by their WHATWG label."),(0,l.kt)("td",{parentName:"tr",align:null},"null, string"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"null"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format/csv/errorThreshold")),(0,l.kt)("td",{parentName:"tr",align:null},"Error Threshold"),(0,l.kt)("td",{parentName:"tr",align:null},"Allows a percentage of errors to be ignored without failing the entire parsing process. When this limit is exceeded, parsing halts."),(0,l.kt)("td",{parentName:"tr",align:null},"integer"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"null"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format/csv/escape")),(0,l.kt)("td",{parentName:"tr",align:null},"Escape Character"),(0,l.kt)("td",{parentName:"tr",align:null},"The escape character, used to escape quotes within fields."),(0,l.kt)("td",{parentName:"tr",align:null},"null, string"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"null"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format/csv/headers")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"Manually specified headers, which can be used in cases where the file itself doesn","'","t contain a header row. If specified, then the parser will assume that the first row is data, not column names, and the column names given here will be used. The column names will be matched with the columns in the file by the order in which they appear here."),(0,l.kt)("td",{parentName:"tr",align:null},"array"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"[]"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format/csv/lineEnding")),(0,l.kt)("td",{parentName:"tr",align:null},"Line Ending"),(0,l.kt)("td",{parentName:"tr",align:null},"The value that terminates a line. Only single-byte values are supported, with the exception of ",'"',"\\","r","\\","n",'"'," (CRLF), which will accept lines terminated by either a carriage return, a newline, or both."),(0,l.kt)("td",{parentName:"tr",align:null},"null, string"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"null"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format/csv/quote")),(0,l.kt)("td",{parentName:"tr",align:null},"Quote Character"),(0,l.kt)("td",{parentName:"tr",align:null},"The character used to quote fields."),(0,l.kt)("td",{parentName:"tr",align:null},"null, string"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"null"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format/json")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"object"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"{}"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/parser/format/w3cExtendedLog")),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null}),(0,l.kt)("td",{parentName:"tr",align:null},"object"),(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"{}"))),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("inlineCode",{parentName:"td"},"/prefix")),(0,l.kt)("td",{parentName:"tr",align:null},"Prefix"),(0,l.kt)("td",{parentName:"tr",align:null},"Prefix within the bucket to capture from."),(0,l.kt)("td",{parentName:"tr",align:null},"string"),(0,l.kt)("td",{parentName:"tr",align:null})),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("strong",{parentName:"td"},(0,l.kt)("inlineCode",{parentName:"strong"},"/region"))),(0,l.kt)("td",{parentName:"tr",align:null},"AWS Region"),(0,l.kt)("td",{parentName:"tr",align:null},"The name of the AWS region where the S3 bucket is located. ",'"',"us-east-1",'"'," is a popular default you can try, if you","'","re unsure what to put here."),(0,l.kt)("td",{parentName:"tr",align:null},"string"),(0,l.kt)("td",{parentName:"tr",align:null},"Required, ",(0,l.kt)("inlineCode",{parentName:"td"},'"us-east-1"'))))),(0,l.kt)("h4",{id:"bindings"},"Bindings"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:null},"Property"),(0,l.kt)("th",{parentName:"tr",align:null},"Title"),(0,l.kt)("th",{parentName:"tr",align:null},"Description"),(0,l.kt)("th",{parentName:"tr",align:null},"Type"),(0,l.kt)("th",{parentName:"tr",align:null},"Required/Default"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("strong",{parentName:"td"},(0,l.kt)("inlineCode",{parentName:"strong"},"/stream"))),(0,l.kt)("td",{parentName:"tr",align:null},"Prefix"),(0,l.kt)("td",{parentName:"tr",align:null},"Path to dataset in the bucket, formatted as ",(0,l.kt)("inlineCode",{parentName:"td"},"bucket-name/prefix-name"),"."),(0,l.kt)("td",{parentName:"tr",align:null},"string"),(0,l.kt)("td",{parentName:"tr",align:null},"Required")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:null},(0,l.kt)("strong",{parentName:"td"},(0,l.kt)("inlineCode",{parentName:"strong"},"/syncMode"))),(0,l.kt)("td",{parentName:"tr",align:null},"Sync mode"),(0,l.kt)("td",{parentName:"tr",align:null},"Connection method. Always set to ",(0,l.kt)("inlineCode",{parentName:"td"},"incremental"),"."),(0,l.kt)("td",{parentName:"tr",align:null},"string"),(0,l.kt)("td",{parentName:"tr",align:null},"Required")))),(0,l.kt)("h3",{id:"sample"},"Sample"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-yaml"},'captures:\n  ${PREFIX}/${CAPTURE_NAME}:\n    endpoint:\n      connector:\n        image: ghcr.io/estuary/source-s3:dev\n        config:\n          bucket: "my-bucket"\n          parser:\n            compression: zip\n            format: csv\n              csv:\n                delimiter: ","\n                encoding: utf-8\n                errorThreshold: 5\n                headers: [ID, username, first_name, last_name]\n                lineEnding: ""\\r"\n                quote: """\n          region: "us-east-1"\n    bindings:\n      - resource:\n          stream: my-bucket/${PREFIX}\n          syncMode: incremental\n        target: ${PREFIX}/${COLLECTION_NAME}\n\n')),(0,l.kt)("p",null,"Your capture definition may be more complex, with additional bindings for different S3 prefixes within the same bucket."),(0,l.kt)("p",null,(0,l.kt)("a",{parentName:"p",href:"/concepts/captures#pull-captures"},"Learn more about capture definitions.")),(0,l.kt)("h3",{id:"advanced-parsing-cloud-storage-data"},"Advanced: Parsing cloud storage data"),(0,l.kt)("p",null,"Cloud storage platforms like S3 can support a wider variety of file types\nthan other data source systems. For each of these file types, Flow must parse\nand translate data into collections with defined fields and JSON schemas."),(0,l.kt)("p",null,"By default, the parser will automatically detect the type and shape of the data in your bucket,\nso you won't need to change the parser configuration for most captures."),(0,l.kt)("p",null,"However, the automatic detection may be incorrect in some cases.\nTo fix or prevent this, you can provide explicit information in the parser configuration,\nwhich is part of the ",(0,l.kt)("a",{parentName:"p",href:"#endpoint"},"endpoint configuration")," for this connector."),(0,l.kt)("p",null,"The parser configuration includes:"),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("strong",{parentName:"p"},"Compression"),": Specify how the bucket contents are compressed.\nIf no compression type is specified, the connector will try to determine the compression type automatically.\nOptions are ",(0,l.kt)("strong",{parentName:"p"},"zip"),", ",(0,l.kt)("strong",{parentName:"p"},"gzip"),", and ",(0,l.kt)("strong",{parentName:"p"},"none"),".")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("strong",{parentName:"p"},"Format"),": Specify the data format, which determines how it will be parsed.\nOptions are:"),(0,l.kt)("ul",{parentName:"li"},(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("strong",{parentName:"p"},"Auto"),": If no format is specified, the connector will try to determine it automatically.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("strong",{parentName:"p"},"Avro"))),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("strong",{parentName:"p"},"CSV"))),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("strong",{parentName:"p"},"JSON"))),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},(0,l.kt)("strong",{parentName:"p"},"W3C Extended Log")),(0,l.kt)("div",{parentName:"li",className:"admonition admonition-info alert alert--info"},(0,l.kt)("div",{parentName:"div",className:"admonition-heading"},(0,l.kt)("h5",{parentName:"div"},(0,l.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,l.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"14",height:"16",viewBox:"0 0 14 16"},(0,l.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"}))),"info")),(0,l.kt)("div",{parentName:"div",className:"admonition-content"},(0,l.kt)("p",{parentName:"div"},"At this time, Flow only supports S3 captures with data of a single file type.\nSupport for multiple file types, which can be configured on a per-binding basis,\nwill be added in the future."),(0,l.kt)("p",{parentName:"div"},"For now, use a prefix in the endpoint configuration to limit the scope of each capture to data of a single file type."))))))),(0,l.kt)("p",null,"Only CSV data requires further configuration. When capturing CSV data, you must specify:"),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("strong",{parentName:"li"},"Delimiter")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("strong",{parentName:"li"},"Encoding")," type, specified by its ",(0,l.kt)("a",{parentName:"li",href:"https://encoding.spec.whatwg.org/#names-and-labels"},"WHATWG label"),"."),(0,l.kt)("li",{parentName:"ul"},"Optionally, an ",(0,l.kt)("strong",{parentName:"li"},"Error threshold"),", as an acceptable percentage of errors."),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("strong",{parentName:"li"},"Escape characters")),(0,l.kt)("li",{parentName:"ul"},"Optionally, a list of column ",(0,l.kt)("strong",{parentName:"li"},"Headers"),", if not already included in the first row of the CSV file."),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("strong",{parentName:"li"},"Line ending")," values"),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("strong",{parentName:"li"},"Quote character"))),(0,l.kt)("p",null,"Descriptions of these properties are included in the ",(0,l.kt)("a",{parentName:"p",href:"#endpoint"},"table above"),"."))}c.isMDXComponent=!0}}]);