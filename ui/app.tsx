import React from "react";

import { ErrorBoundary } from "react-error-boundary";
import {
  ChevronDoubleRightIcon,
  ChevronUpIcon,
  CubeIcon,
  DocumentIcon,
} from "@heroicons/react/24/outline";

import { ArchiveToolbar } from "./components/archive-toolbar/archive-toolbar";
import { Crash } from "./components/crash/crash";
import { Titlebar } from "./components/titlebar/titlebar";
import { TopNav } from "./components/topnav/topnav";

export function AppImpl() {
  return (
    <>
      <Titlebar />
      <TopNav />
      <ArchiveToolbar />
      <div className="breadcrumbs" id="breadcrumbs">
        <CubeIcon width={12} height={12}></CubeIcon>
        <button>/</button>
        <ChevronDoubleRightIcon width={12} height={12}></ChevronDoubleRightIcon>
        <button>data</button>
      </div>
      <div className="file-table-wrap">
        <table className="file-table">
          <thead>
            <tr>
              <th className="check-cell">
                <input
                  id="select-all"
                  type="checkbox"
                  aria-label="Select all visible files"
                />
              </th>
              <th>
                <button>
                  Name{" "}
                  <ChevronUpIcon
                    width={12}
                    height={12}
                    style={{ display: "inline-block" }}
                  ></ChevronUpIcon>
                </button>
              </th>
              <th>Size</th>
              <th>Source archive</th>
              <th className="versions-heading">Versions</th>
            </tr>
          </thead>
          <tbody id="file-rows">
            <tr className="" aria-selected="false">
              <td className="check-cell">
                <input type="checkbox" aria-label="Select pangya.iff" />
              </td>
              <td>
                <button className="file-name-button data">
                  <DocumentIcon width={16} height={16} />
                  pangya.iff
                </button>
              </td>
              <td>13.1 KB</td>
              <td>
                <span className="source-tag">projectg.pak</span>
              </td>
              <td>
                <span>—</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </>
  );
}

export function App() {
  const [errorInfo, setErrorInfo] = React.useState<React.ErrorInfo>();
  return (
    <ErrorBoundary
      fallbackRender={(props) => (
        <Crash componentName="Pangfiles" errorInfo={errorInfo} {...props} />
      )}
      onError={(_, errorInfo) => setErrorInfo(errorInfo)}
    >
      <AppImpl />
    </ErrorBoundary>
  );
}
