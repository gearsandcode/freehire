import{t as h,f as x,g as k,h as w,i as _,j as S}from"./template-tT34KPOZ.js";import{t as D}from"./index-5Po2arxZ.js";import{s as P}from"./snippet-B2J6bAY9.js";import{c as T,a as V,s as A}from"./cn-fA2G5DTu.js";import{p as B}from"./props-Cf8R3aP3.js";import"./branches-D8WOO6Hk.js";const j=D({base:"flex items-start gap-3 rounded-lg border p-4 text-sm",variants:{variant:{default:"border-border bg-card text-foreground",destructive:"border-destructive/30 bg-destructive/5 text-destructive",brand:"border-brand/30 bg-brand-muted text-brand-strong"}},defaultVariants:{variant:"default"}});var q=_('<div role="alert"><!></div>');function v(f,e){w(e,!0);let g=B(e,"variant",3,"default");var s=q(),b=S(s);P(b,()=>e.children),h(y=>A(s,1,y),[()=>T(V(j({variant:g()}),e.class))]),x(f,s),k()}v.__docgen={data:[{name:"variant",visibility:"public",keywords:[],kind:"let",type:{kind:"union",type:[{kind:"const",type:"string",value:"default",text:'"default"'},{kind:"const",type:"string",value:"destructive",text:'"destructive"'},{kind:"const",type:"string",value:"brand",text:'"brand"'}],text:'"default" | "destructive" | "brand"'},static:!1,readonly:!1,defaultValue:'"default"'},{name:"class",visibility:"public",keywords:[],kind:"let",type:{kind:"type",type:"string",text:"string"},static:!1,readonly:!1},{name:"children",visibility:"public",keywords:[{name:"required",description:""}],kind:"let",type:{kind:"function",text:"Snippet<[]>"},static:!1,readonly:!1}],name:"alert.svelte"};const H={title:"Primitives/Alert",component:v,tags:["autodocs"],argTypes:{variant:{control:"select",options:["default","destructive","brand"]}}},t={args:{variant:"default",children:"This is an informational alert."}},a={args:{variant:"destructive",children:"Something went wrong."}},r={args:{variant:"brand",children:"Profile is complete!"}};var n,i,d;t.parameters={...t.parameters,docs:{...(n=t.parameters)==null?void 0:n.docs,source:{originalSource:`{
  args: {
    variant: 'default',
    children: 'This is an informational alert.'
  }
}`,...(d=(i=t.parameters)==null?void 0:i.docs)==null?void 0:d.source}}};var o,l,c;a.parameters={...a.parameters,docs:{...(o=a.parameters)==null?void 0:o.docs,source:{originalSource:`{
  args: {
    variant: 'destructive',
    children: 'Something went wrong.'
  }
}`,...(c=(l=a.parameters)==null?void 0:l.docs)==null?void 0:c.source}}};var u,p,m;r.parameters={...r.parameters,docs:{...(u=r.parameters)==null?void 0:u.docs,source:{originalSource:`{
  args: {
    variant: 'brand',
    children: 'Profile is complete!'
  }
}`,...(m=(p=r.parameters)==null?void 0:p.docs)==null?void 0:m.source}}};const I=["Default","Destructive","Brand"];export{r as Brand,t as Default,a as Destructive,I as __namedExportsOrder,H as default};
