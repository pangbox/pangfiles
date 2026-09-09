import React from "react";

import "./index.css";

interface PushButtonProps extends React.HTMLAttributes<HTMLButtonElement> {
  variant?: "primary" | undefined;
}

export function PushButton({
  variant,
  children,
  ...props
}: React.PropsWithChildren<PushButtonProps>): React.ReactNode {
  const classes = ["push-button", variant].filter(Boolean).join(" ");
  return (
    <button className={classes} {...props}>
      {children}
    </button>
  );
}
