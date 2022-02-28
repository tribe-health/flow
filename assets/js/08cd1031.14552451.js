"use strict";(self.webpackChunksite=self.webpackChunksite||[]).push([[327],{3905:function(e,t,n){n.d(t,{Zo:function(){return u},kt:function(){return d}});var a=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function c(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},o=Object.keys(e);for(a=0;a<o.length;a++)n=o[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(a=0;a<o.length;a++)n=o[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var l=a.createContext({}),s=function(e){var t=a.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},u=function(e){var t=s(e.components);return a.createElement(l.Provider,{value:t},e.children)},p={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},m=a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,o=e.originalType,l=e.parentName,u=c(e,["components","mdxType","originalType","parentName"]),m=s(n),d=r,h=m["".concat(l,".").concat(d)]||m[d]||p[d]||o;return n?a.createElement(h,i(i({ref:t},u),{},{components:n})):a.createElement(h,i({ref:t},u))}));function d(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var o=n.length,i=new Array(o);i[0]=m;var c={};for(var l in t)hasOwnProperty.call(t,l)&&(c[l]=t[l]);c.originalType=e,c.mdxType="string"==typeof e?e:r,i[1]=c;for(var s=2;s<o;s++)i[s]=n[s];return a.createElement.apply(null,i)}return a.createElement.apply(null,n)}m.displayName="MDXCreateElement"},2671:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return c},contentTitle:function(){return l},metadata:function(){return s},toc:function(){return u},default:function(){return m}});var a=n(7462),r=n(3366),o=(n(7294),n(3905)),i=["components"],c={sidebar_position:2},l="Apache Kafka",s={unversionedId:"reference/Connectors/capture-connectors/apache-kafka",id:"reference/Connectors/capture-connectors/apache-kafka",title:"Apache Kafka",description:"This connector captures streaming data from Apache Kafka topics.",source:"@site/docs/reference/Connectors/capture-connectors/apache-kafka.md",sourceDirName:"reference/Connectors/capture-connectors",slug:"/reference/Connectors/capture-connectors/apache-kafka",permalink:"/reference/Connectors/capture-connectors/apache-kafka",editUrl:"https://github.com/estuary/flow/edit/master/site/docs/reference/Connectors/capture-connectors/apache-kafka.md",tags:[],version:"current",sidebarPosition:2,frontMatter:{sidebar_position:2},sidebar:"tutorialSidebar",previous:{title:"Amazon Kinesis",permalink:"/reference/Connectors/capture-connectors/amazon-kinesis"},next:{title:"MySQL",permalink:"/reference/Connectors/capture-connectors/MySQL"}},u=[{value:"Prerequisites",id:"prerequisites",children:[{value:"Authentication and connection security",id:"authentication-and-connection-security",children:[],level:3}],level:2},{value:"Configuration",id:"configuration",children:[{value:"Values",id:"values",children:[],level:3},{value:"Sample",id:"sample",children:[],level:3}],level:2}],p={toc:u};function m(e){var t=e.components,n=(0,r.Z)(e,i);return(0,o.kt)("wrapper",(0,a.Z)({},p,n,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"apache-kafka"},"Apache Kafka"),(0,o.kt)("p",null,"This connector captures streaming data from Apache Kafka topics."),(0,o.kt)("p",null,(0,o.kt)("inlineCode",{parentName:"p"},"ghcr.io/estuary/source-kafka:dev")," provides the latest connector image when using the Flow GitOps environment. You can also follow the link in your browser to see past image versions."),(0,o.kt)("h2",{id:"prerequisites"},"Prerequisites"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"A Kafka cluster with:",(0,o.kt)("ul",{parentName:"li"},(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("a",{parentName:"li",href:"https://kafka.apache.org/documentation/#producerconfigs_bootstrap.servers"},"bootstrap.servers")," configured so that clients may connect via the desired host and port"),(0,o.kt)("li",{parentName:"ul"},"An authentication mechanism of choice set up (highly recommended for production environments)"),(0,o.kt)("li",{parentName:"ul"},"Connection security enabled with TLS (highly recommended for production environments)")))),(0,o.kt)("h3",{id:"authentication-and-connection-security"},"Authentication and connection security"),(0,o.kt)("p",null,"Neither authentication nor connection security are enabled by default in your Kafka cluster, but both are important considerations.\nSimilarly, Flow's Kafka connectors do not strictly require authentication or connection security mechanisms.\nYou may choose to omit them for local development and testing; however, both are strongly encouraged for production environments."),(0,o.kt)("p",null,"A wide ",(0,o.kt)("a",{parentName:"p",href:"https://kafka.apache.org/documentation/#security_overview"},"variety of authentication methods")," is available in Kafka clusters.\nFlow supports SASL/SCRAM-SHA-256, SASL/SCRAM-SHA-512, and SASL/PLAIN. Behavior using other authentication methods is not guaranteed.\nWhen authentication details are not provided, the client connection will attempt to use PLAINTEXT (insecure) protocol."),(0,o.kt)("p",null,"If you don't already have authentication enabled on your cluster, Estuary recommends either of listed ",(0,o.kt)("a",{parentName:"p",href:"https://kafka.apache.org/documentation/#security_sasl_scram"},"SASL/SCRAM")," methods.\nWith SCRAM, you set up a username and password, making it analogous to the traditional authentication mechanisms\nyou use in other applications."),(0,o.kt)("p",null,'For connection security, Estuary recommends that you enable TLS encryption for your SASL mechanism of choice,\nas well as all other components of your cluster.\nNote that because TLS replaced now-deprecated SSL encryption, Kafka still uses the acronym "SSL" to refer to TLS encryption.\nSee ',(0,o.kt)("a",{parentName:"p",href:"https://docs.confluent.io/platform/current/kafka/authentication_ssl.html"},"Confluent's documentation")," for details."),(0,o.kt)("div",{className:"admonition admonition-info alert alert--info"},(0,o.kt)("div",{parentName:"div",className:"admonition-heading"},(0,o.kt)("h5",{parentName:"div"},(0,o.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,o.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"14",height:"16",viewBox:"0 0 14 16"},(0,o.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"}))),"Beta")),(0,o.kt)("div",{parentName:"div",className:"admonition-content"},(0,o.kt)("p",{parentName:"div"},"TLS encryption is currently the only supported connection security mechanism for this connector.\nOther connection security methods may be enabled in the future."))),(0,o.kt)("h2",{id:"configuration"},"Configuration"),(0,o.kt)("p",null,"There are various ways to configure and implement connectors. See ",(0,o.kt)("a",{parentName:"p",href:"/concepts/connectors#using-connectors"},"connectors")," to learn more about these methods. The values and code sample below provide configuration details specific to the Apache Kafka source connector."),(0,o.kt)("h3",{id:"values"},"Values"),(0,o.kt)("table",null,(0,o.kt)("thead",{parentName:"table"},(0,o.kt)("tr",{parentName:"thead"},(0,o.kt)("th",{parentName:"tr",align:null},"Value"),(0,o.kt)("th",{parentName:"tr",align:null},"Name"),(0,o.kt)("th",{parentName:"tr",align:null},"Description"),(0,o.kt)("th",{parentName:"tr",align:null},"Type"),(0,o.kt)("th",{parentName:"tr",align:null},"Required/Default"))),(0,o.kt)("tbody",{parentName:"table"},(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},(0,o.kt)("inlineCode",{parentName:"td"},"bootstrap_servers")),(0,o.kt)("td",{parentName:"tr",align:null},"Bootstrap servers"),(0,o.kt)("td",{parentName:"tr",align:null},"The initial servers in the Kafka cluster to connect to. The Kafka client will be informed of the rest of the cluster nodes by connecting to one of these nodes."),(0,o.kt)("td",{parentName:"tr",align:null},"array"),(0,o.kt)("td",{parentName:"tr",align:null},"Required")),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},(0,o.kt)("inlineCode",{parentName:"td"},"tls")),(0,o.kt)("td",{parentName:"tr",align:null},"TLS"),(0,o.kt)("td",{parentName:"tr",align:null},"TLS connection settings"),(0,o.kt)("td",{parentName:"tr",align:null},"string"),(0,o.kt)("td",{parentName:"tr",align:null},(0,o.kt)("inlineCode",{parentName:"td"},'"system_certificates"'))),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},(0,o.kt)("inlineCode",{parentName:"td"},"authentication")),(0,o.kt)("td",{parentName:"tr",align:null},"Authentication"),(0,o.kt)("td",{parentName:"tr",align:null},"Connection details used to authenticate a client connection to Kafka via SASL"),(0,o.kt)("td",{parentName:"tr",align:null},"null, object"),(0,o.kt)("td",{parentName:"tr",align:null})),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},(0,o.kt)("inlineCode",{parentName:"td"},"authentication/mechanism")),(0,o.kt)("td",{parentName:"tr",align:null},"Mechanism"),(0,o.kt)("td",{parentName:"tr",align:null},"SASL mechanism describing how to exchange and authenticate client servers"),(0,o.kt)("td",{parentName:"tr",align:null},"string"),(0,o.kt)("td",{parentName:"tr",align:null})),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},(0,o.kt)("inlineCode",{parentName:"td"},"authentication/password")),(0,o.kt)("td",{parentName:"tr",align:null},"Password"),(0,o.kt)("td",{parentName:"tr",align:null},"Password, if applicable for the authentication mechanism chosen"),(0,o.kt)("td",{parentName:"tr",align:null},"string"),(0,o.kt)("td",{parentName:"tr",align:null})),(0,o.kt)("tr",{parentName:"tbody"},(0,o.kt)("td",{parentName:"tr",align:null},(0,o.kt)("inlineCode",{parentName:"td"},"authentication/username")),(0,o.kt)("td",{parentName:"tr",align:null},"Username"),(0,o.kt)("td",{parentName:"tr",align:null},"Username, if applicable for the authentication mechanism chosen"),(0,o.kt)("td",{parentName:"tr",align:null},"string"),(0,o.kt)("td",{parentName:"tr",align:null})))),(0,o.kt)("h3",{id:"sample"},"Sample"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},"captures:\n  ${TENANT}/${CAPTURE_NAME}:\n    endpoint:\n      connector:\n        image: ghcr.io/estuary/source-kafka:dev\n        config:\n            bootstrap_servers: [localhost:9093]\n            tls: system_certificates\n            authentication:\n                mechanism: SCRAM-SHA-512\n                username: bruce.wayne\n                password: definitely-not-batman\n    bindings:\n      - resource:\n           stream: ${STREAM_NAME}\n           syncMode: incremental\n        target: ${TENANT}/${COLLECTION_NAME}\n")))}m.isMDXComponent=!0}}]);