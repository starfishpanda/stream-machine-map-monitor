/// <reference types="vite/client" />
/// <reference types="vite/types/importMeta.d.ts" />

// CSS modules declaration
declare module '*.css';

// Extend the existing Vite ImportMetaEnv interface
interface ImportMetaEnv extends Readonly<Record<string, string | boolean | undefined>> {
  readonly VITE_GOOGLE_MAPS_API_KEY: string;
  // You can add other VITE_ prefixed vars here
}

// Extend ImportMeta
interface ImportMeta {
  readonly env: ImportMetaEnv;
}
