import React from "react";

import "./index.css";

interface ToolButtonProps extends React.HTMLAttributes<HTMLButtonElement> {}

export function ToolButton({
  children,
  ...props
}: React.PropsWithChildren<ToolButtonProps>): React.ReactNode {
  return (
    <button className="tool-button" {...props}>
      {children}
    </button>
  );
}
