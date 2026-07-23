import{t as S,f as V,g as C,h as D,i as j,j as w}from"./template-tT34KPOZ.js";import{t as P}from"./index-5Po2arxZ.js";import{s as A}from"./snippet-B2J6bAY9.js";import{c as B,a as R,s as q}from"./cn-fA2G5DTu.js";import{p as E}from"./props-Cf8R3aP3.js";import"./branches-D8WOO6Hk.js";const O=P({base:"inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium transition-colors",variants:{variant:{default:"border-border bg-muted text-foreground",primary:"border-transparent bg-primary text-primary-foreground",secondary:"border-transparent bg-secondary text-secondary-foreground",brand:"border-transparent bg-brand-muted text-brand-strong",destructive:"border-destructive/20 bg-destructive/10 text-destructive"}},defaultVariants:{variant:"default"}});var T=j("<span><!></span>");function b(h,e){D(e,!0);let x=E(e,"variant",3,"default");var s=T(),k=w(s);A(k,()=>e.children),S(_=>q(s,1,_),[()=>B(R(O({variant:x()}),e.class))]),V(h,s),C()}b.__docgen={data:[{name:"variant",visibility:"public",keywords:[],kind:"let",type:{kind:"union",type:[{kind:"const",type:"string",value:"default",text:'"default"'},{kind:"const",type:"string",value:"primary",text:'"primary"'},{kind:"const",type:"string",value:"secondary",text:'"secondary"'},{kind:"const",type:"string",value:"brand",text:'"brand"'},{kind:"const",type:"string",value:"destructive",text:'"destructive"'}],text:'"default" | "primary" | "secondary" | "brand" | "destructive"'},static:!1,readonly:!1,defaultValue:'"default"'},{name:"class",visibility:"public",keywords:[],kind:"let",type:{kind:"type",type:"string",text:"string"},static:!1,readonly:!1},{name:"children",visibility:"public",keywords:[{name:"required",description:""}],kind:"let",type:{kind:"function",text:"Snippet<[]>"},static:!1,readonly:!1}],name:"chip.svelte"};const K={title:"Primitives/Chip",component:b,tags:["autodocs"],argTypes:{variant:{control:"select",options:["default","primary","secondary","brand","destructive"]}}},r={args:{variant:"default",children:"Chip"}},t={args:{variant:"primary",children:"Active"}},a={args:{variant:"brand",children:"Verified"}},n={args:{variant:"destructive",children:"Rejected"}};var i,d,o;r.parameters={...r.parameters,docs:{...(i=r.parameters)==null?void 0:i.docs,source:{originalSource:`{
  args: {
    variant: 'default',
    children: 'Chip'
  }
}`,...(o=(d=r.parameters)==null?void 0:d.docs)==null?void 0:o.source}}};var c,p,l;t.parameters={...t.parameters,docs:{...(c=t.parameters)==null?void 0:c.docs,source:{originalSource:`{
  args: {
    variant: 'primary',
    children: 'Active'
  }
}`,...(l=(p=t.parameters)==null?void 0:p.docs)==null?void 0:l.source}}};var u,m,v;a.parameters={...a.parameters,docs:{...(u=a.parameters)==null?void 0:u.docs,source:{originalSource:`{
  args: {
    variant: 'brand',
    children: 'Verified'
  }
}`,...(v=(m=a.parameters)==null?void 0:m.docs)==null?void 0:v.source}}};var f,y,g;n.parameters={...n.parameters,docs:{...(f=n.parameters)==null?void 0:f.docs,source:{originalSource:`{
  args: {
    variant: 'destructive',
    children: 'Rejected'
  }
}`,...(g=(y=n.parameters)==null?void 0:y.docs)==null?void 0:g.source}}};const L=["Default","Primary","Brand","Destructive"];export{a as Brand,r as Default,n as Destructive,t as Primary,L as __namedExportsOrder,K as default};
