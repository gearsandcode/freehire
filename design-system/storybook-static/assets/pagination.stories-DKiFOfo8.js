import{t as G,f as y,g as H,h as I,n as a,i as J,j as v,o as d,u as p,k as K,l as Q}from"./template-tT34KPOZ.js";import{s as R}from"./render-zwLsI_mv.js";import{s as T}from"./snippet-B2J6bAY9.js";import{i as U}from"./if-QPuj07Ph.js";import{d as x,a as W}from"./utils-Bqw0_cdZ.js";import{c as X,a as Y,s as Z}from"./cn-fA2G5DTu.js";import{p as P,u as k}from"./props-Cf8R3aP3.js";import"./branches-D8WOO6Hk.js";var $=J('<nav aria-label="Pagination"><button type="button" class="inline-flex h-9 items-center justify-center rounded-md border border-border px-3 text-sm transition-colors hover:bg-accent disabled:pointer-events-none disabled:opacity-50" aria-label="Previous page">Previous</button> <span class="text-sm text-muted-foreground"> </span> <button type="button" class="inline-flex h-9 items-center justify-center rounded-md border border-border px-3 text-sm transition-colors hover:bg-accent disabled:pointer-events-none disabled:opacity-50" aria-label="Next page">Next</button> <!></nav>');function V(q,e){I(e,!0);let t=P(e,"page",15,1),E=P(e,"perPage",3,20),m=p(()=>Math.max(1,Math.ceil(e.total/E()))),g=p(()=>t()>1),u=p(()=>t()<a(m));function O(){a(g)&&k(t,-1)}function z(){a(u)&&k(t)}var o=$(),l=v(o),b=d(l,2),A=v(b),c=d(b,2),B=d(c,2);{var C=r=>{var f=K(),F=Q(f);T(F,()=>e.children),y(r,f)};U(B,r=>{e.children&&r(C)})}G(r=>{Z(o,1,r),l.disabled=!a(g),R(A,`Page ${t()??""} of ${a(m)??""}`),c.disabled=!a(u)},[()=>X(Y("flex items-center gap-2",e.class))]),x("click",l,O),x("click",c,z),y(q,o),H()}W(["click"]);V.__docgen={data:[{name:"page",visibility:"public",keywords:[],kind:"let",type:{kind:"type",type:"number",text:"number"},static:!1,readonly:!1,defaultValue:"..."},{name:"total",visibility:"public",keywords:[{name:"required",description:""}],kind:"let",type:{kind:"type",type:"number",text:"number"},static:!1,readonly:!1},{name:"perPage",visibility:"public",keywords:[],kind:"let",type:{kind:"type",type:"number",text:"number"},static:!1,readonly:!1,defaultValue:"20"},{name:"class",visibility:"public",keywords:[],kind:"let",type:{kind:"type",type:"string",text:"string"},static:!1,readonly:!1},{name:"children",visibility:"public",keywords:[],kind:"let",type:{kind:"function",text:"Snippet<[]>"},static:!1,readonly:!1}],name:"pagination.svelte"};const le={title:"Primitives/Pagination",component:V,tags:["autodocs"]},s={args:{page:1,total:250,perPage:20}},n={args:{page:13,total:250,perPage:20}},i={args:{page:1,total:5,perPage:20}};var h,_,S;s.parameters={...s.parameters,docs:{...(h=s.parameters)==null?void 0:h.docs,source:{originalSource:`{
  args: {
    page: 1,
    total: 250,
    perPage: 20
  }
}`,...(S=(_=s.parameters)==null?void 0:_.docs)==null?void 0:S.source}}};var w,j,N;n.parameters={...n.parameters,docs:{...(w=n.parameters)==null?void 0:w.docs,source:{originalSource:`{
  args: {
    page: 13,
    total: 250,
    perPage: 20
  }
}`,...(N=(j=n.parameters)==null?void 0:j.docs)==null?void 0:N.source}}};var D,L,M;i.parameters={...i.parameters,docs:{...(D=i.parameters)==null?void 0:D.docs,source:{originalSource:`{
  args: {
    page: 1,
    total: 5,
    perPage: 20
  }
}`,...(M=(L=i.parameters)==null?void 0:L.docs)==null?void 0:M.source}}};const ce=["Default","LastPage","SinglePage"];export{s as Default,n as LastPage,i as SinglePage,ce as __namedExportsOrder,le as default};
