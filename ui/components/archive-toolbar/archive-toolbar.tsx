import React from "react";
import { ArrowDownTrayIcon, DocumentPlusIcon, ServerIcon } from "@heroicons/react/24/outline";

import "./index.css";

import { ToolButton } from "../../widgets/tool-button/tool-button";

export function ArchiveToolbar(): React.ReactNode {
  return (
    <div className="workspace-toolbar" aria-label="Archive actions">
      <ToolButton>
        <DocumentPlusIcon width={16} height={16} />
        Add archives…
      </ToolButton>
      <span className="toolbar-separator"></span>
      <ToolButton>
        <ArrowDownTrayIcon width={16} height={16} />
        Extract…
      </ToolButton>
      <ToolButton>
        <ServerIcon width={16} height={16} />
        Mount…
      </ToolButton>
      <label className="region-picker">
        <span>Region</span>
        <select id="region" aria-label="Archive region">
          <option value="">Auto-detect</option>
          <option value="us">United States</option>
          <option value="jp">Japan</option>
          <option value="th">Thailand</option>
          <option value="eu">Europe</option>
          <option value="id">Indonesia</option>
          <option value="kr">Korea</option>
        </select>
      </label>
    </div>
  );
}
