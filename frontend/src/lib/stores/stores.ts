import { writable } from 'svelte/store';

export const selectedTabID = writable(null);

export const errors = writable([]);