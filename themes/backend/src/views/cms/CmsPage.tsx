import React from 'react'
import { Header } from '@/components/header/Header'
import { Footer } from '@/components/footer/Footer'

const CmsPage: React.FC = () => (
  <>
    <Header />
    <div className="page-wrapper">
      <div className="container-fluid">
        <div className="page-header d-print-none">
          <div className="row align-items-center">
            <div className="col">
              <div className="page-pretitle">
                Overview
              </div>
              <h2 className="page-title">
                Fluid layout
              </h2>
            </div>
            <div className="col-auto ms-auto d-print-none">
              <div className="btn-list">
                  <span className="d-none d-sm-inline">
                    <a href="#" className="btn btn-white">
                      New view
                    </a>
                  </span>
                <a href="#" className="btn btn-primary d-none d-sm-inline-block" data-bs-toggle="modal"
                   data-bs-target="#modal-report">
                  <svg xmlns="http://www.w3.org/2000/svg" className="icon" width="24" height="24" viewBox="0 0 24 24"
                       stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round"
                       stroke-linejoin="round">
                    <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                    <line x1="12" y1="5" x2="12" y2="19"></line>
                    <line x1="5" y1="12" x2="19" y2="12"></line>
                  </svg>
                  Create new report
                </a>
                <a href="#" className="btn btn-primary d-sm-none btn-icon" data-bs-toggle="modal"
                   data-bs-target="#modal-report" aria-label="Create new report">
                  <svg xmlns="http://www.w3.org/2000/svg" className="icon" width="24" height="24" viewBox="0 0 24 24"
                       stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round"
                       stroke-linejoin="round">
                    <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                    <line x1="12" y1="5" x2="12" y2="19"></line>
                    <line x1="5" y1="12" x2="19" y2="12"></line>
                  </svg>
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div className="page-body">
        <div className="container-fluid">
          <div className="row row-deck row-cards">
            <div className="col-sm-6 col-lg-3">
              <div className="card">
                <div className="card-body">
                  <div className="d-flex align-items-center">
                    <div className="subheader">Sales</div>
                    <div className="ms-auto lh-1">
                      <div className="dropdown">
                        <a className="dropdown-toggle text-muted" href="#" data-bs-toggle="dropdown"
                           aria-haspopup="true" aria-expanded="false">Last 7 days</a>
                        <div className="dropdown-menu dropdown-menu-end">
                          <a className="dropdown-item active" href="#">Last 7 days</a>
                          <a className="dropdown-item" href="#">Last 30 days</a>
                          <a className="dropdown-item" href="#">Last 3 months</a>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div className="h1 mb-3">75%</div>
                  <div className="d-flex mb-2">
                    <div>Conversion rate</div>
                    <div className="ms-auto">
                        <span className="text-green d-inline-flex align-items-center lh-1">
                          7%
                          <svg xmlns="http://www.w3.org/2000/svg" className="icon ms-1" width="24" height="24"
                               viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none"
                               stroke-linecap="round" stroke-linejoin="round"><path stroke="none" d="M0 0h24v24H0z"
                                                                                    fill="none"></path><polyline
                            points="3 17 9 11 13 15 21 7"></polyline><polyline
                            points="14 7 21 7 21 14"></polyline></svg>
                        </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <Footer />
    </div>
  </>
)

export default CmsPage
