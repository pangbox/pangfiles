import React from "react";

import "./index.css";
import { CloseButton } from "../../widgets/close-button/close-button";

export function Titlebar(): React.ReactNode {
  return (
    <header className="titlebar">
      <a className="program-name">Pangfiles</a>
      <span className="document-title">C:\PangYa\ProjectG*.pak</span>
      <div className="window-tools">
        <CloseButton />
      </div>
    </header>
  );
}
