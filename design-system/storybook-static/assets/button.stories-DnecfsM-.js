import{k as Q,l as R,f as m,g as U,h as W,i as H,j as v}from"./template-tT34KPOZ.js";import{t as X}from"./index-5Po2arxZ.js";import{s as f}from"./snippet-B2J6bAY9.js";import{i as Y}from"./if-QPuj07Ph.js";import{b as h}from"./attributes-Bp59jm6V.js";import{p as b,r as Z}from"./props-Cf8R3aP3.js";import{a as k}from"./cn-fA2G5DTu.js";import"./branches-D8WOO6Hk.js";import"./utils-Bqw0_cdZ.js";const x=X({base:"inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50",variants:{variant:{primary:"bg-brand text-brand-foreground hover:opacity-90",secondary:"bg-secondary text-secondary-foreground hover:bg-accent",outline:"border border-border bg-background hover:bg-accent hover:text-accent-foreground",ghost:"hover:bg-accent hover:text-accent-foreground"},size:{sm:"h-8 px-3",md:"h-9 px-4",lg:"h-11 px-6",icon:"size-8"}},defaultVariants:{variant:"secondary",size:"md"}});var $=new Set(["$$slots","$$events","$$legacy","variant","size","class","href","children"]),ee=H("<a><!></a>"),te=H("<button><!></button>");function I(J,e){W(e,!0);let u=b(e,"variant",3,"secondary"),p=b(e,"size",3,"md"),y=Z(e,$);var g=Q(),K=R(g);{var M=r=>{var t=ee();h(t,d=>({href:e.href,class:d,...y}),[()=>k(x({variant:u(),size:p()}),e.class)]);var l=v(t);f(l,()=>e.children),m(r,t)},N=r=>{var t=te();h(t,d=>({type:"button",class:d,...y}),[()=>k(x({variant:u(),size:p()}),e.class)]);var l=v(t);f(l,()=>e.children),m(r,t)};Y(K,r=>{e.href?r(M):r(N,-1)})}m(J,g),U()}I.__docgen={data:[{name:"variant",visibility:"public",keywords:[],kind:"let",type:{kind:"union",type:[{kind:"const",type:"string",value:"secondary",text:'"secondary"'},{kind:"const",type:"string",value:"primary",text:'"primary"'},{kind:"const",type:"string",value:"outline",text:'"outline"'},{kind:"const",type:"string",value:"ghost",text:'"ghost"'}],text:'"secondary" | "primary" | "outline" | "ghost"'},static:!1,readonly:!1,defaultValue:'"secondary"'},{name:"size",visibility:"public",keywords:[],kind:"let",type:{kind:"union",type:[{kind:"const",type:"string",value:"md",text:'"md"'},{kind:"const",type:"string",value:"sm",text:'"sm"'},{kind:"const",type:"string",value:"lg",text:'"lg"'},{kind:"const",type:"string",value:"icon",text:'"icon"'}],text:'"md" | "sm" | "lg" | "icon"'},static:!1,readonly:!1,defaultValue:'"md"'},{name:"class",visibility:"public",keywords:[],kind:"let",type:{kind:"union",type:[{kind:"type",type:"string",text:"string"},{kind:"type",type:"intersection",text:"string & ClassArray"},{kind:"type",type:"intersection",text:"string & ClassDictionary"}],text:"string | string & ClassArray | string & ClassDictionary"},static:!1,readonly:!1},{name:"href",visibility:"public",keywords:[],kind:"let",type:{kind:"type",type:"string",text:"string"},static:!1,readonly:!1},{name:"children",visibility:"public",keywords:[{name:"required",description:""}],kind:"let",type:{kind:"function",text:"Snippet<[]>"},static:!1,readonly:!1},{name:"type",visibility:"public",keywords:[],kind:"let",type:{kind:"union",type:[{kind:"const",type:"string",value:"submit",text:'"submit"'},{kind:"const",type:"string",value:"reset",text:'"reset"'},{kind:"const",type:"string",value:"button",text:'"button"'}],text:'"submit" | "reset" | "button"'},static:!1,readonly:!1}],name:"button.svelte"};const me={title:"Primitives/Button",component:I,tags:["autodocs"],argTypes:{variant:{control:"select",options:["primary","secondary","outline","ghost"]},size:{control:"select",options:["sm","md","lg","icon"]}}},a={args:{variant:"secondary",size:"md",children:"Click me"}},n={args:{variant:"primary",size:"md",children:"Primary"}},s={args:{variant:"outline",size:"md",children:"Outline"}},i={args:{variant:"ghost",size:"md",children:"Ghost"}},o={args:{variant:"secondary",size:"sm",children:"Small"}},c={args:{variant:"primary",size:"lg",children:"Large"}};var z,_,S;a.parameters={...a.parameters,docs:{...(z=a.parameters)==null?void 0:z.docs,source:{originalSource:`{
  args: {
    variant: 'secondary',
    size: 'md',
    children: 'Click me'
  }
}`,...(S=(_=a.parameters)==null?void 0:_.docs)==null?void 0:S.source}}};var w,C,O;n.parameters={...n.parameters,docs:{...(w=n.parameters)==null?void 0:w.docs,source:{originalSource:`{
  args: {
    variant: 'primary',
    size: 'md',
    children: 'Primary'
  }
}`,...(O=(C=n.parameters)==null?void 0:C.docs)==null?void 0:O.source}}};var P,D,G;s.parameters={...s.parameters,docs:{...(P=s.parameters)==null?void 0:P.docs,source:{originalSource:`{
  args: {
    variant: 'outline',
    size: 'md',
    children: 'Outline'
  }
}`,...(G=(D=s.parameters)==null?void 0:D.docs)==null?void 0:G.source}}};var L,V,j;i.parameters={...i.parameters,docs:{...(L=i.parameters)==null?void 0:L.docs,source:{originalSource:`{
  args: {
    variant: 'ghost',
    size: 'md',
    children: 'Ghost'
  }
}`,...(j=(V=i.parameters)==null?void 0:V.docs)==null?void 0:j.source}}};var A,B,q;o.parameters={...o.parameters,docs:{...(A=o.parameters)==null?void 0:A.docs,source:{originalSource:`{
  args: {
    variant: 'secondary',
    size: 'sm',
    children: 'Small'
  }
}`,...(q=(B=o.parameters)==null?void 0:B.docs)==null?void 0:q.source}}};var E,T,F;c.parameters={...c.parameters,docs:{...(E=c.parameters)==null?void 0:E.docs,source:{originalSource:`{
  args: {
    variant: 'primary',
    size: 'lg',
    children: 'Large'
  }
}`,...(F=(T=c.parameters)==null?void 0:T.docs)==null?void 0:F.source}}};const ue=["Default","Primary","Outline","Ghost","Small","Large"];export{a as Default,i as Ghost,c as Large,s as Outline,n as Primary,o as Small,ue as __namedExportsOrder,me as default};
