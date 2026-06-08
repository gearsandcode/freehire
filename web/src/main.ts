import { mount } from 'svelte';
import App from './App.svelte';
import './app.css';
import { initTheme } from '$lib/theme.svelte';

// Apply the saved (or OS) theme once before mounting so there is no flash of
// the wrong palette.
initTheme();

const target = document.getElementById('app');
if (!target) throw new Error('#app element not found in index.html');

const app = mount(App, { target });

export default app;
