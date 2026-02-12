"use client";

import { useEffect, useRef } from "react";
import { cn } from "@/lib/utils";

type AvatarDisplayProps = {
  readonly seed: string;
  readonly size?: number;
  readonly className?: string;
};

export function AvatarDisplay({
  seed,
  size = 80,
  className,
}: AvatarDisplayProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    const hash = simpleHash(seed);
    const colors = generateColors(hash);

    canvas.width = size;
    canvas.height = size;

    const cellSize = size / 5;
    for (let x = 0; x < 3; x++) {
      for (let y = 0; y < 5; y++) {
        const idx = x * 5 + y;
        if ((hash >> idx) & 1) {
          ctx.fillStyle = colors[idx % colors.length];
          ctx.fillRect(x * cellSize, y * cellSize, cellSize, cellSize);
          ctx.fillRect((4 - x) * cellSize, y * cellSize, cellSize, cellSize);
        }
      }
    }
  }, [seed, size]);

  return (
    <canvas
      ref={canvasRef}
      width={size}
      height={size}
      className={cn("rounded-lg", className)}
    />
  );
}

function simpleHash(str: string): number {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i);
    hash = (hash << 5) - hash + char;
    hash = hash & hash;
  }
  return Math.abs(hash);
}

function generateColors(hash: number): readonly string[] {
  const hue = hash % 360;
  return [
    `hsl(${hue}, 70%, 50%)`,
    `hsl(${(hue + 40) % 360}, 70%, 60%)`,
    `hsl(${(hue + 80) % 360}, 60%, 50%)`,
  ] as const;
}
