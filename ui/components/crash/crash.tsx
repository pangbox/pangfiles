import { ExclamationTriangleIcon } from "@heroicons/react/24/solid";
import React from "react";

import { Modal } from "../../widgets/modal/modal";
import { PushButton } from "../../widgets/push-button/push-button";

import "./index.css";

interface CrashProps {
  error?: unknown | undefined;
  errorInfo?: React.ErrorInfo | undefined;
  resetErrorBoundary?: (() => void) | undefined;
  componentName?: string | undefined;
}

export function Crash({
  error,
  errorInfo,
  resetErrorBoundary,
  componentName = "Component",
}: React.PropsWithChildren<CrashProps>): React.ReactNode {
  const [isOpen, setIsOpen] = React.useState(true);

  const reset = React.useCallback(() => {
    resetErrorBoundary?.();
    setIsOpen(false);
  }, [resetErrorBoundary])

  const newIssueURL =
    "https://github.com/pangbox/pangfiles/issues/new" +
    `?title=Crash+in+${encodeURIComponent(componentName)}` +
    `&body=Component+error:+${encodeURIComponent(String(error))}`;
  return (
    <>
    <div className="crash-placeholder">
        <span>{componentName} has crashed.</span>
        <PushButton onClick={() => setIsOpen(true)}>Show Details</PushButton>
    </div>
    <Modal
      header={`${componentName} has crashed`}
      footer={
        resetErrorBoundary ?
        <PushButton variant="primary" onClick={reset}>Reset</PushButton> : null
      }
      isOpen={isOpen}
      onClose={() => setIsOpen(false)}
    >
      <div className="crash-content">
        <div className="crash-icon-wrapper">
          <ExclamationTriangleIcon width={32} height={32} />
        </div>

        <p>
          An error occurred during component rendering. This is a bug.{" "}
          <a href={newIssueURL}>Open a bug report on GitHub?</a>
        </p>

        {error !== undefined && error instanceof Error && (
          <div className="crash-error-details">
            <div className="crash-error-name">{error.name}</div>
            <div className="crash-error-message">{error.message}</div>
          </div>
        )}

        {errorInfo && (
          <details className="crash-stack-details">
            <summary className="crash-stack-summary">Render Stack</summary>
            <pre className="crash-stack-trace">
              {String(errorInfo.componentStack).trim()}
            </pre>
          </details>
        )}
      </div>
      </Modal>
    </>
  );
}
