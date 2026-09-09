import React from "react";

import "./index.css";
import { XMarkIcon } from "@heroicons/react/24/solid";

interface CloseButtonProps extends React.HTMLAttributes<HTMLButtonElement> {}

export function CloseButton(props: CloseButtonProps): React.ReactNode {
  return (
    <button className="close-button" {...props}>
      <XMarkIcon></XMarkIcon>
    </button>
  );
}
