(function(e){function t(t){for(var r,o,u=t[0],i=t[1],s=t[2],l=0,f=[];l<u.length;l++)o=u[l],Object.prototype.hasOwnProperty.call(a,o)&&a[o]&&f.push(a[o][0]),a[o]=0;for(r in i)Object.prototype.hasOwnProperty.call(i,r)&&(e[r]=i[r]);d&&d(t);while(f.length)f.shift()();return c.push.apply(c,s||[]),n()}function n(){for(var e,t=0;t<c.length;t++){for(var n=c[t],r=!0,o=1;o<n.length;o++){var i=n[o];0!==a[i]&&(r=!1)}r&&(c.splice(t--,1),e=u(u.s=n[0]))}return e}var r={},a={app:0},c=[];function o(e){return u.p+"js/"+({about:"about"}[e]||e)+"."+{about:"0cd2f7b7"}[e]+".js"}function u(t){if(r[t])return r[t].exports;var n=r[t]={i:t,l:!1,exports:{}};return e[t].call(n.exports,n,n.exports,u),n.l=!0,n.exports}u.e=function(e){var t=[],n=a[e];if(0!==n)if(n)t.push(n[2]);else{var r=new Promise((function(t,r){n=a[e]=[t,r]}));t.push(n[2]=r);var c,i=document.createElement("script");i.charset="utf-8",i.timeout=120,u.nc&&i.setAttribute("nonce",u.nc),i.src=o(e);var s=new Error;c=function(t){i.onerror=i.onload=null,clearTimeout(l);var n=a[e];if(0!==n){if(n){var r=t&&("load"===t.type?"missing":t.type),c=t&&t.target&&t.target.src;s.message="Loading chunk "+e+" failed.\n("+r+": "+c+")",s.name="ChunkLoadError",s.type=r,s.request=c,n[1](s)}a[e]=void 0}};var l=setTimeout((function(){c({type:"timeout",target:i})}),12e4);i.onerror=i.onload=c,document.head.appendChild(i)}return Promise.all(t)},u.m=e,u.c=r,u.d=function(e,t,n){u.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:n})},u.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},u.t=function(e,t){if(1&t&&(e=u(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var n=Object.create(null);if(u.r(n),Object.defineProperty(n,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var r in e)u.d(n,r,function(t){return e[t]}.bind(null,r));return n},u.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return u.d(t,"a",t),t},u.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},u.p="/",u.oe=function(e){throw console.error(e),e};var i=window["webpackJsonp"]=window["webpackJsonp"]||[],s=i.push.bind(i);i.push=t,i=i.slice();for(var l=0;l<i.length;l++)t(i[l]);var d=s;c.push([0,"chunk-vendors"]),n()})({0:function(e,t,n){e.exports=n("56d7")},1:function(e,t){},10:function(e,t){},2:function(e,t){},3:function(e,t){},4:function(e,t){},5:function(e,t){},"56d7":function(e,t,n){"use strict";n.r(t);n("e260"),n("e6cf"),n("cca6"),n("a79d");var r=n("7a23");function a(e,t,n,a,c,o){var u=Object(r["F"])("router-view");return Object(r["y"])(),Object(r["f"])(u)}var c=n("5530"),o=n("dd2f"),u=n("365c"),i=n("31a8"),s=n("ad17"),l=n("a012"),d=n("957f"),f=n("d2f5"),b=n("b070"),p=n("d703"),O=n("f4db"),j=n("f685"),g=n("fd32"),h=n.n(g),m={setup:function(){h.a.init({liffId:"1656247924-eX5ZOvN0"}).then((function(){h.a.isLoggedIn()||h.a.login()}));var e=new o["a"]((function(e,t){return e.setContext((function(e){var t=e.headers,n=void 0===t?{}:t;return{headers:Object(c["a"])(Object(c["a"])({},n),{},{Authorization:h.a.isLoggedIn()?h.a.getAccessToken():null})}})),t(e)})),t=new u["a"]({link:Object(i["a"])((function(e){var t=e.query,n=Object(b["e"])(t);return"OperationDefinition"===n.kind&&"subscription"===n.operation}),new f["a"]({uri:"wss://tinychats.herokuapp.com/graphql",options:{reconnect:!0}}),Object(s["a"])([e,Object(p["a"])({useGETForHashedQueries:!0,sha256:O["sha256"]}),Object(l["a"])({uri:"https://tinychats.herokuapp.com/graphql"})])),cache:new d["a"]});Object(r["A"])(j["a"],t)}};n("6d6f");m.render=a;var v=m,y=(n("d3b7"),n("3ca3"),n("ddb0"),n("6c02")),x=(n("b0c0"),n("498a"),Object(r["N"])("data-v-828441dc"));Object(r["B"])("data-v-828441dc");var w=Object(r["g"])("Loading..."),M=Object(r["h"])("div",{id:"message-end",style:{"margin-bottom":"60px"}},null,-1),k=Object(r["g"])("Send");Object(r["z"])();var _,S,L,U,C=x((function(e,t,n,a,c,o){var u=Object(r["F"])("van-loading"),i=Object(r["F"])("van-image"),s=Object(r["F"])("van-badge"),l=Object(r["F"])("van-cell"),d=Object(r["F"])("van-button"),f=Object(r["F"])("van-field");return Object(r["y"])(),Object(r["f"])(r["a"],null,[a.currentUserLoading||a.listMessagesLoading?(Object(r["y"])(),Object(r["f"])(u,{key:0,style:{"text-align":"center","margin-top":"10px"}},{default:x((function(){return[w]})),_:1})):(Object(r["y"])(!0),Object(r["f"])(r["a"],{key:1},Object(r["E"])(a.messages,(function(e,t){return Object(r["y"])(),Object(r["f"])(l,{key:e.id},{title:x((function(){return[Object(r["h"])(s,{content:e.user.name,color:"#1989fa",style:{width:"max-content",left:"0px",top:"3px"}},{default:x((function(){return[Object(r["h"])(i,{src:e.user.avatarUrl,width:"30px",height:"30px",round:""},null,8,["src"])]})),_:2},1032,["content"])]})),default:x((function(){return[Object(r["g"])(Object(r["I"])(e.text),1)]})),_:2},1024)})),128)),(Object(r["y"])(!0),Object(r["f"])(r["a"],null,Object(r["E"])(a.messagesCreated,(function(e,t){return Object(r["y"])(),Object(r["f"])(l,{key:e.id},{title:x((function(){return[Object(r["h"])(s,{content:e.user.name,color:"#1989fa",style:{width:"max-content",left:"0px",top:"3px"}},{default:x((function(){return[Object(r["h"])(i,{src:e.user.avatarUrl,width:"30px",height:"30px",round:""},null,8,["src"])]})),_:2},1032,["content"])]})),default:x((function(){return[Object(r["g"])(Object(r["I"])(e.text),1)]})),_:2},1024)})),128)),M,Object(r["h"])(f,{class:"fixedbutton",modelValue:a.createMessageState,"onUpdate:modelValue":t[1]||(t[1]=function(e){return a.createMessageState=e}),size:"small",placeholder:"please input message"},{button:x((function(){return[Object(r["h"])(d,{size:"small",icon:"comment-o",loading:a.createMessageLoading,disabled:""===a.createMessageState.trim(),onClick:a.createMessage},{default:x((function(){return[k]})),_:1},8,["loading","disabled","onClick"])]})),_:1},8,["modelValue"])],64)})),P=n("8785"),A=n("5184"),F=n("f672"),q=Object(A["a"])(_||(_=Object(P["a"])(["\n  query currentUser {\n    currentUser {\n      id\n      name\n      avatarUrl\n    }\n  }\n"]))),E=Object(A["a"])(S||(S=Object(P["a"])(["\n  query listMessages {\n    messages {\n      id\n      text\n      createdAt\n      user {\n        id\n        name\n        avatarUrl\n      }\n    }\n  }\n"]))),I=Object(A["a"])(L||(L=Object(P["a"])(["\n  mutation createMessage($text: String!) {\n    createMessage(input: { text: $text }) {\n      id\n      text\n      createdAt\n      user {\n        id\n        name\n        avatarUrl\n      }\n    }\n  }\n"]))),T=Object(A["a"])(U||(U=Object(P["a"])(["\n  subscription onMessageCreated {\n    messageCreated {\n      id\n      text\n      createdAt\n      user {\n        id\n        name\n        avatarUrl\n      }\n    }\n  }\n"]))),z={name:"Home",setup:function(){var e=Object(j["c"])(q),t=e.loading,n=Object(j["c"])(E),a=n.result,c=n.loading,o=n.onResult,u=Object(j["d"])(a,[],(function(e){return e.messages})),i=Object(r["D"])([]),s=Object(r["D"])(""),l=Object(j["b"])(I,(function(){return{variables:{text:s.value}}})),d=l.mutate,f=l.loading,b=l.onDone;b((function(){return s.value=""}));var p=Object(j["e"])(T),O=p.result;return Object(r["L"])(O,(function(e){i.value.push(JSON.parse(JSON.stringify(e.messageCreated))),Object(F["a"])("#message-end")}),{lazy:!0}),o((function(){return Object(F["a"])(window.innerHeight)})),{currentUserLoading:t,listMessagesLoading:c,messages:u,messagesCreated:i,createMessageState:s,createMessage:d,createMessageLoading:f}}};n("cebf");z.render=C,z.__scopeId="data-v-828441dc";var D=z,H=[{path:"/",name:"Home",component:D},{path:"/about",name:"About",component:function(){return n.e("about").then(n.bind(null,"f820"))}}],J=Object(y["a"])({history:Object(y["b"])("/"),routes:H}),N=J,V=n("b970"),$=(n("157a"),function(e){return e.use(V["a"])}),B=Object(r["e"])(v);$(B),B.use(N).mount("#app")},6:function(e,t){},6026:function(e,t,n){},"6d6f":function(e,t,n){"use strict";n("f15b")},7:function(e,t){},8:function(e,t){},9:function(e,t){},cebf:function(e,t,n){"use strict";n("6026")},f15b:function(e,t,n){}});
//# sourceMappingURL=app.2574eeb8.js.map