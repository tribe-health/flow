"use strict";(self.webpackChunksite=self.webpackChunksite||[]).push([[6882],{3905:function(e,t,a){a.d(t,{Zo:function(){return d},kt:function(){return h}});var n=a(7294);function r(e,t,a){return t in e?Object.defineProperty(e,t,{value:a,enumerable:!0,configurable:!0,writable:!0}):e[t]=a,e}function o(e,t){var a=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),a.push.apply(a,n)}return a}function s(e){for(var t=1;t<arguments.length;t++){var a=null!=arguments[t]?arguments[t]:{};t%2?o(Object(a),!0).forEach((function(t){r(e,t,a[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(a)):o(Object(a)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(a,t))}))}return e}function i(e,t){if(null==e)return{};var a,n,r=function(e,t){if(null==e)return{};var a,n,r={},o=Object.keys(e);for(n=0;n<o.length;n++)a=o[n],t.indexOf(a)>=0||(r[a]=e[a]);return r}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(n=0;n<o.length;n++)a=o[n],t.indexOf(a)>=0||Object.prototype.propertyIsEnumerable.call(e,a)&&(r[a]=e[a])}return r}var c=n.createContext({}),l=function(e){var t=n.useContext(c),a=t;return e&&(a="function"==typeof e?e(t):s(s({},t),e)),a},d=function(e){var t=l(e.components);return n.createElement(c.Provider,{value:t},e.children)},p={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},u=n.forwardRef((function(e,t){var a=e.components,r=e.mdxType,o=e.originalType,c=e.parentName,d=i(e,["components","mdxType","originalType","parentName"]),u=l(a),h=r,m=u["".concat(c,".").concat(h)]||u[h]||p[h]||o;return a?n.createElement(m,s(s({ref:t},d),{},{components:a})):n.createElement(m,s({ref:t},d))}));function h(e,t){var a=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var o=a.length,s=new Array(o);s[0]=u;var i={};for(var c in t)hasOwnProperty.call(t,c)&&(i[c]=t[c]);i.originalType=e,i.mdxType="string"==typeof e?e:r,s[1]=i;for(var l=2;l<o;l++)s[l]=a[l];return n.createElement.apply(null,s)}return n.createElement.apply(null,a)}u.displayName="MDXCreateElement"},4594:function(e,t,a){a.r(t),a.d(t,{frontMatter:function(){return i},contentTitle:function(){return c},metadata:function(){return l},toc:function(){return d},default:function(){return u}});var n=a(7462),r=a(3366),o=(a(7294),a(3905)),s=["components"],i={},c="Task shards",l={unversionedId:"concepts/advanced/shards",id:"concepts/advanced/shards",title:"Task shards",description:"Task shards are an advanced concept of Flow.",source:"@site/docs/concepts/advanced/shards.md",sourceDirName:"concepts/advanced",slug:"/concepts/advanced/shards",permalink:"/concepts/advanced/shards",editUrl:"https://github.com/estuary/flow/edit/master/site/docs/concepts/advanced/shards.md",tags:[],version:"current",frontMatter:{},sidebar:"tutorialSidebar",previous:{title:"Projections",permalink:"/concepts/advanced/projections"},next:{title:"Create a simple data flow",permalink:"/guides/create-dataflow"}},d=[{value:"Shard splits",id:"shard-splits",children:[],level:2},{value:"Recovery logs",id:"recovery-logs",children:[],level:2}],p={toc:d};function u(e){var t=e.components,a=(0,r.Z)(e,s);return(0,o.kt)("wrapper",(0,n.Z)({},p,a,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"task-shards"},"Task shards"),(0,o.kt)("div",{className:"admonition admonition-tip alert alert--success"},(0,o.kt)("div",{parentName:"div",className:"admonition-heading"},(0,o.kt)("h5",{parentName:"div"},(0,o.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,o.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"12",height:"16",viewBox:"0 0 12 16"},(0,o.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M6.5 0C3.48 0 1 2.19 1 5c0 .92.55 2.25 1 3 1.34 2.25 1.78 2.78 2 4v1h5v-1c.22-1.22.66-1.75 2-4 .45-.75 1-2.08 1-3 0-2.81-2.48-5-5.5-5zm3.64 7.48c-.25.44-.47.8-.67 1.11-.86 1.41-1.25 2.06-1.45 3.23-.02.05-.02.11-.02.17H5c0-.06 0-.13-.02-.17-.2-1.17-.59-1.83-1.45-3.23-.2-.31-.42-.67-.67-1.11C2.44 6.78 2 5.65 2 5c0-2.2 2.02-4 4.5-4 1.22 0 2.36.42 3.22 1.19C10.55 2.94 11 3.94 11 5c0 .66-.44 1.78-.86 2.48zM4 14h5c-.23 1.14-1.3 2-2.5 2s-2.27-.86-2.5-2z"}))),"tip")),(0,o.kt)("div",{parentName:"div",className:"admonition-content"},(0,o.kt)("p",{parentName:"div"},"Task shards are an advanced concept of Flow.\nYou can use Flow without knowing the details of shards,\nbut this section may help you better understand how Flow works."))),(0,o.kt)("p",null,"Catalog ",(0,o.kt)("a",{parentName:"p",href:"/concepts/#tasks"},"tasks")," \u2014 captures, derivations, and materializations \u2014\nare executed by one or more task ",(0,o.kt)("strong",{parentName:"p"},"shards"),"."),(0,o.kt)("p",null,"Shards are a fault-tolerant and stateful unit of execution for a catalog task,\nwhich the Flow runtime assigns and runs on a scalable pool of compute resources.\nA single task can have many shards,\nwhich allow the task to scale across many machines to\nachieve more throughput and parallelism."),(0,o.kt)("p",null,"Shards are part of the Gazette project.\n",(0,o.kt)("a",{parentName:"p",href:"https://gazette.readthedocs.io/en/latest/consumers-concepts.html#shards"},"See Gazette's Shard concepts page for details"),"."),(0,o.kt)("h2",{id:"shard-splits"},"Shard splits"),(0,o.kt)("p",null,"When a task is first created, it is initialized with a single shard.\nLater and as required, a shard can be split into two shards.\nOnce initiated, the split may require up to a few minutes to complete,\nbut it doesn't require downtime and the selected shard continues\nto run until the split occurs."),(0,o.kt)("p",null,"This process can be repeated as needed until your required throughput is achieved."),(0,o.kt)("div",{className:"admonition admonition-caution alert alert--warning"},(0,o.kt)("div",{parentName:"div",className:"admonition-heading"},(0,o.kt)("h5",{parentName:"div"},(0,o.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,o.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"16",height:"16",viewBox:"0 0 16 16"},(0,o.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M8.893 1.5c-.183-.31-.52-.5-.887-.5s-.703.19-.886.5L.138 13.499a.98.98 0 0 0 0 1.001c.193.31.53.501.886.501h13.964c.367 0 .704-.19.877-.5a1.03 1.03 0 0 0 .01-1.002L8.893 1.5zm.133 11.497H6.987v-2.003h2.039v2.003zm0-3.004H6.987V5.987h2.039v4.006z"}))),"TODO")),(0,o.kt)("div",{parentName:"div",className:"admonition-content"},(0,o.kt)("p",{parentName:"div"},"This section is incomplete.\nSee ",(0,o.kt)("inlineCode",{parentName:"p"},"flowctl shards split --help")," for further details."))),(0,o.kt)("h2",{id:"recovery-logs"},"Recovery logs"),(0,o.kt)("div",{className:"admonition admonition-info alert alert--info"},(0,o.kt)("div",{parentName:"div",className:"admonition-heading"},(0,o.kt)("h5",{parentName:"div"},(0,o.kt)("span",{parentName:"h5",className:"admonition-icon"},(0,o.kt)("svg",{parentName:"span",xmlns:"http://www.w3.org/2000/svg",width:"14",height:"16",viewBox:"0 0 14 16"},(0,o.kt)("path",{parentName:"svg",fillRule:"evenodd",d:"M7 2.3c3.14 0 5.7 2.56 5.7 5.7s-2.56 5.7-5.7 5.7A5.71 5.71 0 0 1 1.3 8c0-3.14 2.56-5.7 5.7-5.7zM7 1C3.14 1 0 4.14 0 8s3.14 7 7 7 7-3.14 7-7-3.14-7-7-7zm1 3H6v5h2V4zm0 6H6v2h2v-2z"}))),"info")),(0,o.kt)("div",{parentName:"div",className:"admonition-content"},(0,o.kt)("p",{parentName:"div"},"Shard stores and associated states are transparent to you, the Flow user.\nThis section is informational only, to provide a sense of how Flow works."))),(0,o.kt)("p",null,"All task shards have associated state, which is managed in the shard's store."),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"Capture tasks must track incremental checkpoints of their endpoint connectors."),(0,o.kt)("li",{parentName:"ul"},"Derivation tasks manage a potentially very large index of registers,\nas well as read checkpoints of sourced collection journals."),(0,o.kt)("li",{parentName:"ul"},"Materialization tasks track incremental checkpoints of their endpoint connectors,\nas well as read checkpoints of sourced collection journals.")),(0,o.kt)("p",null,"Shard stores use\n",(0,o.kt)("a",{parentName:"p",href:"https://gazette.readthedocs.io/en/latest/consumers-concepts.html#recovery-logs"},"recovery logs"),"\nto replicate updates and implement transaction semantics."),(0,o.kt)("p",null,"Recovery logs are regular ",(0,o.kt)("a",{parentName:"p",href:"/concepts/advanced/journals"},"journals"),",\nbut hold binary data and are not intended for direct use.\nHowever, they can hold your user data.\nRecovery logs of ",(0,o.kt)("a",{parentName:"p",href:"/concepts/derivations"},"derivations")," hold your derivation register values."),(0,o.kt)("p",null,"Recovery logs are stored in your cloud storage bucket,\nand must have a configured ",(0,o.kt)("a",{parentName:"p",href:"/concepts/storage-mappings#recovery-logs"},"storage mapping"),"."))}u.isMDXComponent=!0}}]);