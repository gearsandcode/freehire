import{t as S,f as _,g as w,h as M,i as O,j as V}from"./template-tT34KPOZ.js";import{t as N}from"./index-5Po2arxZ.js";import{s as j}from"./snippet-B2J6bAY9.js";import{c as q,a as E,s as P}from"./cn-fA2G5DTu.js";import{p as T}from"./props-Cf8R3aP3.js";import"./branches-D8WOO6Hk.js";const z=N({base:"inline-flex items-center rounded-md border px-2 py-0.5 text-xs font-medium",variants:{variant:{secondary:"border-transparent bg-secondary text-secondary-foreground",outline:"border-border text-muted-foreground",brand:"border-transparent bg-brand-muted text-brand-strong",missing:"border-destructive/15 bg-destructive/5 text-destructive/75"}},defaultVariants:{variant:"secondary"}});var A=O("<span><!></span>");function f(x,e){M(e,!0);let h=T(e,"variant",3,"secondary");var s=A(),k=V(s);j(k,()=>e.children),S(B=>P(s,1,B),[()=>q(E(z({variant:h()}),e.class))]),_(x,s),w()}f.__docgen={data:[{name:"variant",visibility:"public",keywords:[],kind:"let",type:{kind:"union",type:[{kind:"const",type:"string",value:"secondary",text:'"secondary"'},{kind:"const",type:"string",value:"outline",text:'"outline"'},{kind:"const",type:"string",value:"brand",text:'"brand"'},{kind:"const",type:"string",value:"missing",text:'"missing"'}],text:'"secondary" | "outline" | "brand" | "missing"'},static:!1,readonly:!1,defaultValue:'"secondary"'},{name:"class",visibility:"public",keywords:[],kind:"let",type:{kind:"type",type:"string",text:"string"},static:!1,readonly:!1},{name:"children",visibility:"public",keywords:[{name:"required",description:""}],kind:"let",type:{kind:"function",text:"Snippet<[]>"},static:!1,readonly:!1}],name:"badge.svelte"};const J={title:"Primitives/Badge",component:f,tags:["autodocs"],argTypes:{variant:{control:"select",options:["secondary","outline","brand","missing"]}}},a={args:{variant:"secondary",children:"Badge"}},r={args:{variant:"outline",children:"Badge"}},n={args:{variant:"brand",children:"New"}},t={args:{variant:"missing",children:"Missing"}};var i,o,d;a.parameters={...a.parameters,docs:{...(i=a.parameters)==null?void 0:i.docs,source:{originalSource:`{
  args: {
    variant: 'secondary',
    children: 'Badge'
  }
}`,...(d=(o=a.parameters)==null?void 0:o.docs)==null?void 0:d.source}}};var c,l,p;r.parameters={...r.parameters,docs:{...(c=r.parameters)==null?void 0:c.docs,source:{originalSource:`{
  args: {
    variant: 'outline',
    children: 'Badge'
  }
}`,...(p=(l=r.parameters)==null?void 0:l.docs)==null?void 0:p.source}}};var m,u,g;n.parameters={...n.parameters,docs:{...(m=n.parameters)==null?void 0:m.docs,source:{originalSource:`{
  args: {
    variant: 'brand',
    children: 'New'
  }
}`,...(g=(u=n.parameters)==null?void 0:u.docs)==null?void 0:g.source}}};var y,v,b;t.parameters={...t.parameters,docs:{...(y=t.parameters)==null?void 0:y.docs,source:{originalSource:`{
  args: {
    variant: 'missing',
    children: 'Missing'
  }
}`,...(b=(v=t.parameters)==null?void 0:v.docs)==null?void 0:b.source}}};const K=["Secondary","Outline","Brand","Missing"];export{n as Brand,t as Missing,r as Outline,a as Secondary,K as __namedExportsOrder,J as default};
