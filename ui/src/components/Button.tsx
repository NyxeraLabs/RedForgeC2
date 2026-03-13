import type { ButtonHTMLAttributes } from "react";

type Props = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: "primary" | "secondary" | "danger";
};

export function Button({ variant = "secondary", className = "", ...props }: Props) {
  const base =
    "inline-flex items-center justify-center rounded px-4 py-2 text-sm font-medium disabled:opacity-50";
  const styles =
    variant === "primary"
      ? "bg-indigo-600 hover:bg-indigo-500"
      : variant === "danger"
        ? "bg-red-600 hover:bg-red-500"
        : "bg-slate-800 hover:bg-slate-700";

  return <button {...props} className={`${base} ${styles} ${className}`} />;
}

