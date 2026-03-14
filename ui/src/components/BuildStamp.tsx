import React from "react";

function buildSha(): string | null {
  const sha = (import.meta as any)?.env?.VITE_BUILD_SHA as string | undefined;
  if (!sha) return null;
  const trimmed = sha.trim();
  return trimmed.length ? trimmed : null;
}

function buildEpoch(): string | null {
  const epoch = (import.meta as any)?.env?.VITE_BUILD_EPOCH as string | undefined;
  if (!epoch) return null;
  const trimmed = epoch.trim();
  return trimmed.length ? trimmed : null;
}

/**
 * BuildStamp renders a tiny UI build identifier (when available) so operators can
 * confirm that a freshly rebuilt container is actually serving new assets.
 *
 * This is visual-only; it does not change any runtime behavior.
 */
export function BuildStamp(props: { className?: string }) {
  const sha = buildSha();
  const epoch = buildEpoch();
  if (!sha && !epoch) return null;
  return (
    <span className={props.className} title="UI build id">
      UI:{sha ?? "dev"}{epoch ? `@${epoch}` : ""}
    </span>
  );
}
