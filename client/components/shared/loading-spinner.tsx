import { cn } from "@/lib/utils";

type LoadingSpinnerProps = {
  readonly className?: string;
  readonly size?: "sm" | "md" | "lg";
};

const sizeClasses = {
  sm: "h-4 w-4 border",
  md: "h-8 w-8 border-2",
  lg: "h-12 w-12 border-2",
} as const;

export function LoadingSpinner({
  className,
  size = "md",
}: LoadingSpinnerProps) {
  return (
    <div
      className={cn(
        "animate-spin rounded-full border-primary border-t-transparent",
        sizeClasses[size],
        className,
      )}
    />
  );
}
