import React from "react";

import "./index.css";
import {
  CircleStackIcon,
  ServerStackIcon,
  Square3Stack3DIcon,
} from "@heroicons/react/24/outline";

export function TopNav(): React.ReactNode {
  return (
    <nav className="topnav" aria-label="Main navigation">
      <button
        className="nav-tab active"
        aria-current="page"
      >
        <Square3Stack3DIcon width={16} height={16} />
        Archives
      </button>
      <button className="nav-tab" aria-current="false">
        <CircleStackIcon width={16} height={16} />
        Update lists
      </button>
      <button className="nav-tab" aria-current="false">
        <ServerStackIcon width={16} height={16} />
        Server
      </button>
    </nav>
  );
}
