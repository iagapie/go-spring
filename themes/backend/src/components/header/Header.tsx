import React from 'react'

import { appName, routes } from '@/utils/constants'

export const Header: React.FC = () => (
  <header className="navbar navbar-expand-md navbar-dark d-print-none">
    <div className="container-fluid">
      <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbar-menu">
        <span className="navbar-toggler-icon" />
      </button>
      <h1 className="navbar-brand navbar-brand-autodark d-none-navbar-horizontal pe-0 pe-md-3">
        <a href={routes.root} className="text-decoration-none">
          {appName}
        </a>
      </h1>
      <div className="navbar-nav flex-row order-md-last">
        <div className="nav-item dropdown">
          <a href="#" className="nav-link d-flex lh-1 text-reset p-0" data-bs-toggle="dropdown"
             aria-label="Open user menu">
            <span className="avatar avatar-sm" />
            <div className="d-none d-xl-block ps-2">
              <div>Igor Agapie</div>
              <div className="mt-1 small text-muted">Developer</div>
            </div>
          </a>
          <div className="dropdown-menu dropdown-menu-end dropdown-menu-arrow">
            <a href="#" className="dropdown-item">Set status</a>
            <a href="#" className="dropdown-item">Profile &amp; account</a>
            <a href="#" className="dropdown-item">Feedback</a>
            <div className="dropdown-divider" />
            <a href="#" className="dropdown-item">Settings</a>
            <a href="#" className="dropdown-item">Logout</a>
          </div>
        </div>
      </div>
      <div className="collapse navbar-collapse" id="navbar-menu">
        <div className="d-flex flex-column flex-md-row flex-fill align-items-stretch align-items-md-center">
          <ul className="navbar-nav">
            <li className="nav-item">
              <a className="nav-link" href={routes.root}>
                <span className="nav-link-icon d-md-none d-lg-inline-block">
                  <svg xmlns="http://www.w3.org/2000/svg" className="icon icon-tabler icon-tabler-dashboard" width="24"
                       height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none"
                       stroke-linecap="round" stroke-linejoin="round">
                   <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                   <circle cx="12" cy="13" r="2" />
                   <line x1="13.45" y1="11.55" x2="15.5" y2="9.5" />
                   <path d="M6.4 20a9 9 0 1 1 11.2 0z" />
                  </svg>
                </span>
                <span className="nav-link-title">
                  Dashboard
                </span>
              </a>
            </li>
            <li className="nav-item dropdown active">
              <a className="nav-link dropdown-toggle show" href="#" data-bs-toggle="dropdown" role="button" aria-expanded="true">
                <span className="nav-link-icon d-md-none d-lg-inline-block">
                  <svg xmlns="http://www.w3.org/2000/svg" className="icon" width="24" height="24"
                       viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round"
                       stroke-linejoin="round">
                    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
                    <rect x="4" y="4" width="6" height="5" rx="2" />
                    <rect x="4" y="13" width="6" height="7" rx="2" />
                    <rect x="14" y="4" width="6" height="7" rx="2" />
                    <rect x="14" y="15" width="6" height="5" rx="2" />
                  </svg>
                </span>
                <span className="nav-link-title">CMS</span>
              </a>
              <div className="dropdown-menu show" data-bs-popper="none">
                <a className="dropdown-item active" href={routes.cms.pages}>
                  Pages
                </a>
                <a className="dropdown-item" href={routes.cms.partials}>
                  Partials
                </a>
                <a className="dropdown-item" href={routes.cms.layouts}>
                  Layouts
                </a>
                <a className="dropdown-item" href={routes.cms.components}>
                  Components
                </a>
                <a className="dropdown-item" href={routes.cms.assets}>
                  Assets
                </a>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </header>
)
