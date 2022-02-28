"use strict";(self.webpackChunksite=self.webpackChunksite||[]).push([[7856],{3905:function(e,t,n){n.d(t,{Zo:function(){return u},kt:function(){return f}});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function c(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var l=r.createContext({}),s=function(e){var t=r.useContext(l),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},u=function(e){var t=s(e.components);return r.createElement(l.Provider,{value:t},e.children)},m={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},p=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,o=e.originalType,l=e.parentName,u=c(e,["components","mdxType","originalType","parentName"]),p=s(n),f=a,k=p["".concat(l,".").concat(f)]||p[f]||m[f]||o;return n?r.createElement(k,i(i({ref:t},u),{},{components:n})):r.createElement(k,i({ref:t},u))}));function f(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=n.length,i=new Array(o);i[0]=p;var c={};for(var l in t)hasOwnProperty.call(t,l)&&(c[l]=t[l]);c.originalType=e,c.mdxType="string"==typeof e?e:a,i[1]=c;for(var s=2;s<o;s++)i[s]=n[s];return r.createElement.apply(null,i)}return r.createElement.apply(null,n)}p.displayName="MDXCreateElement"},2473:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return c},contentTitle:function(){return l},metadata:function(){return s},toc:function(){return u},default:function(){return p}});var r=n(7462),a=n(3366),o=(n(7294),n(3905)),i=["components"],c={},l="Materialization connectors",s={unversionedId:"reference/Connectors/materialization-connectors/README",id:"reference/Connectors/materialization-connectors/README",title:"Materialization connectors",description:"Estuary's available materialization connectors are listed in this section. Each connector has a unique configuration you must follow in your catalog specification; these will be linked below the connector name.",source:"@site/docs/reference/Connectors/materialization-connectors/README.md",sourceDirName:"reference/Connectors/materialization-connectors",slug:"/reference/Connectors/materialization-connectors/",permalink:"/reference/Connectors/materialization-connectors/",editUrl:"https://github.com/estuary/flow/edit/master/site/docs/reference/Connectors/materialization-connectors/README.md",tags:[],version:"current",frontMatter:{},sidebar:"tutorialSidebar",previous:{title:"PostgreSQL",permalink:"/reference/Connectors/capture-connectors/PostgreSQL"},next:{title:"Google BigQuery",permalink:"/reference/Connectors/materialization-connectors/BigQuery"}},u=[{value:"Available materialization connectors",id:"available-materialization-connectors",children:[],level:2}],m={toc:u};function p(e){var t=e.components,n=(0,a.Z)(e,i);return(0,o.kt)("wrapper",(0,r.Z)({},m,n,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"materialization-connectors"},"Materialization connectors"),(0,o.kt)("p",null,"Estuary's available materialization connectors are listed in this section. Each connector has a unique configuration you must follow in your ",(0,o.kt)("a",{parentName:"p",href:"/concepts/#specifications"},"catalog specification"),"; these will be linked below the connector name."),(0,o.kt)("div",{className:"admonition admonition-info alert alert--info"},(0,o.kt)("div",{parentName:"div",className:"admonition-heading"},(0,o.kt)("h5",{parentName:"div"},(0,o.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,o.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"14",height:"16",viewBox:"0 0 14 16"},(0,o.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"}))),"Beta")),(0,o.kt)("div",{parentName:"div",className:"admonition-content"},(0,o.kt)("p",{parentName:"div"},"More configurations coming to the docs soon. ",(0,o.kt)("a",{parentName:"p",href:"mailto:info@estuary.dev"},"Contact the team")," for more information on missing connectors."))),(0,o.kt)("p",null,"Also listed are links to the most recent Docker image, which you'll need for certain configuration methods."),(0,o.kt)("p",null,"Estuary is actively developing new connectors, so check back regularly for the latest additions. We\u2019re prioritizing the development of high-scale technological systems, as well as client needs."),(0,o.kt)("h2",{id:"available-materialization-connectors"},"Available materialization connectors"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"Apache Parquet",(0,o.kt)("ul",{parentName:"li"},(0,o.kt)("li",{parentName:"ul"},"Configuration"),(0,o.kt)("li",{parentName:"ul"},"Package \u2014 ghcr.io/estuary/materialize-s3-parquet:dev"))),(0,o.kt)("li",{parentName:"ul"},"Elasticsearch",(0,o.kt)("ul",{parentName:"li"},(0,o.kt)("li",{parentName:"ul"},"Configuration"),(0,o.kt)("li",{parentName:"ul"},"Package \u2014 ghcr.io/estuary/materialize-elasticsearch:dev"))),(0,o.kt)("li",{parentName:"ul"},"Google BigQuery",(0,o.kt)("ul",{parentName:"li"},(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("a",{parentName:"li",href:"/reference/Connectors/materialization-connectors/BigQuery"},"Configuration")),(0,o.kt)("li",{parentName:"ul"},"Package \u2014 ghcr.io/estuary/materialize-bigquery:dev"))),(0,o.kt)("li",{parentName:"ul"},"PostgreSQL",(0,o.kt)("ul",{parentName:"li"},(0,o.kt)("li",{parentName:"ul"},"Configuration"),(0,o.kt)("li",{parentName:"ul"},"Package \u2014 ghcr.io/estuary/materialize-postgres:dev"))),(0,o.kt)("li",{parentName:"ul"},"Rockset",(0,o.kt)("ul",{parentName:"li"},(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("a",{parentName:"li",href:"/reference/Connectors/materialization-connectors/Rockset"},"Configuration")),(0,o.kt)("li",{parentName:"ul"},"Package \u2014 ghcr.io/estuary/materialize-rockset:dev"))),(0,o.kt)("li",{parentName:"ul"},"Snowflake",(0,o.kt)("ul",{parentName:"li"},(0,o.kt)("li",{parentName:"ul"},"Configuration"),(0,o.kt)("li",{parentName:"ul"},"Package \u2014 ghcr.io/estuary/materialize-snowflake:dev"))),(0,o.kt)("li",{parentName:"ul"},"Webhook",(0,o.kt)("ul",{parentName:"li"},(0,o.kt)("li",{parentName:"ul"},"Configuration"),(0,o.kt)("li",{parentName:"ul"},"Package \u2014 ghcr.io/estuary/materialize-webhook:dev")))))}p.isMDXComponent=!0}}]);