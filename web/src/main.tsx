import './styles/global.scss'

import { createRoot } from 'react-dom/client'

import { App } from './app'

declare global {
  interface Window {
    readonly app: HTMLElement
  }
}

createRoot(window.app).render(<App />)
