import React from "react";
import { createPortal } from "react-dom";

import "./index.css";
import { CloseButton } from "../close-button/close-button";

interface ModalProps {
  isOpen?: boolean | undefined;
  onClose?: (() => void) | undefined;
  header?: React.ReactNode | undefined;
  footer?: React.ReactNode | undefined;
}

export function Modal({
  isOpen = true,
  onClose,
  header,
  footer,
  children,
}: React.PropsWithChildren<ModalProps>): React.ReactNode {
  const dialogRef = React.useRef<HTMLDialogElement>(null);

  React.useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;
    if (isOpen) {
      if (!dialog.open) dialog.showModal();
    } else {
      if (dialog.open) dialog.close();
    }
  }, [isOpen]);

  React.useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;
    const handleCancel = (event: Event) => {
      event.preventDefault();
      onClose?.();
    };
    dialog.addEventListener("cancel", handleCancel);
    return () => dialog.removeEventListener("cancel", handleCancel);
  }, [onClose]);

  return createPortal(
    <dialog className="modal" ref={dialogRef}>
      <div className="modal-content">
        <div className="modal-header">
          {header}
          <CloseButton aria-label="Close dialog" onClick={onClose} />
        </div>
        <div className="modal-body">{children}</div>
        <div className="modal-footer">{footer}</div>
      </div>
    </dialog>,
    document.body,
  );
}
