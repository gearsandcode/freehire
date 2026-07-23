import{k as V,l as q,f as m,g as B,h as F,t as g,n as p,i as C,u as f,j as G}from"./template-tT34KPOZ.js";import{s as I}from"./render-zwLsI_mv.js";import{i as K}from"./if-QPuj07Ph.js";import{s as l,a as Q}from"./attributes-Bp59jm6V.js";import{c as y,a as v,s as k}from"./cn-fA2G5DTu.js";import{p as R}from"./props-Cf8R3aP3.js";import"./utils-Bqw0_cdZ.js";import"./branches-D8WOO6Hk.js";var X=C("<img/>"),Y=C("<div> </div>");function L(W,a){F(a,!0);let c=R(a,"size",3,"md");const d={sm:"size-8 text-xs",md:"size-10 text-sm",lg:"size-12 text-base"};function E(t){let e=0;for(let s=0;s<t.length;s++)e=e*31+t.charCodeAt(s)|0;return Math.abs(e)%360}let H=f(()=>a.name?a.name.split(" ").slice(0,2).map(t=>{var e;return((e=t[0])==null?void 0:e.toUpperCase())??""}).join(""):"?"),M=f(()=>a.name?`hsl(${E(a.name)} 45% 90%)`:"hsl(0 0% 90%)");var u=V(),O=q(u);{var P=t=>{var e=X();g(s=>{l(e,"src",a.src),l(e,"alt",a.name??"avatar"),k(e,1,s)},[()=>y(v("rounded-full object-cover",d[c()],a.class))]),m(t,e)},T=t=>{var e=Y(),s=G(e);g(U=>{k(e,1,U),Q(e,`background-color: ${p(M)??""}`),l(e,"aria-label",a.name??"avatar"),I(s,p(H))},[()=>y(v("flex items-center justify-center rounded-full font-medium text-foreground",d[c()],a.class))]),m(t,e)};K(O,t=>{a.src?t(P):t(T,-1)})}m(W,u),B()}L.__docgen={data:[{name:"name",visibility:"public",keywords:[],kind:"let",type:{kind:"type",type:"string",text:"string"},static:!1,readonly:!1},{name:"src",visibility:"public",keywords:[],kind:"let",type:{kind:"type",type:"string",text:"string"},static:!1,readonly:!1},{name:"size",visibility:"public",keywords:[],kind:"let",type:{kind:"union",type:[{kind:"const",type:"string",value:"sm",text:'"sm"'},{kind:"const",type:"string",value:"md",text:'"md"'},{kind:"const",type:"string",value:"lg",text:'"lg"'}],text:'"sm" | "md" | "lg"'},static:!1,readonly:!1,defaultValue:'"md"'},{name:"class",visibility:"public",keywords:[],kind:"let",type:{kind:"type",type:"string",text:"string"},static:!1,readonly:!1}],name:"avatar.svelte"};const ne={title:"Primitives/Avatar",component:L,tags:["autodocs"],argTypes:{size:{control:"select",options:["sm","md","lg"]}}},r={args:{name:"Jane Doe",size:"md"}},i={args:{name:"John Smith",size:"sm"}},n={args:{name:"Alice Wonderland",size:"lg"}},o={args:{size:"md"}};var x,h,b;r.parameters={...r.parameters,docs:{...(x=r.parameters)==null?void 0:x.docs,source:{originalSource:`{
  args: {
    name: 'Jane Doe',
    size: 'md'
  }
}`,...(b=(h=r.parameters)==null?void 0:h.docs)==null?void 0:b.source}}};var z,_,S;i.parameters={...i.parameters,docs:{...(z=i.parameters)==null?void 0:z.docs,source:{originalSource:`{
  args: {
    name: 'John Smith',
    size: 'sm'
  }
}`,...(S=(_=i.parameters)==null?void 0:_.docs)==null?void 0:S.source}}};var A,j,w;n.parameters={...n.parameters,docs:{...(A=n.parameters)==null?void 0:A.docs,source:{originalSource:`{
  args: {
    name: 'Alice Wonderland',
    size: 'lg'
  }
}`,...(w=(j=n.parameters)==null?void 0:j.docs)==null?void 0:w.source}}};var D,J,N;o.parameters={...o.parameters,docs:{...(D=o.parameters)==null?void 0:D.docs,source:{originalSource:`{
  args: {
    size: 'md'
  }
}`,...(N=(J=o.parameters)==null?void 0:J.docs)==null?void 0:N.source}}};const oe=["Default","Small","Large","NoName"];export{r as Default,n as Large,o as NoName,i as Small,oe as __namedExportsOrder,ne as default};
