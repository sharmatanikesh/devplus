import { Zap } from "lucide-react";

interface LoadingPulseProps {
  className?: string;
  size?: number;
}

export function LoadingPulse({ className = "", size = 24 }: LoadingPulseProps) {
  return (
    <div className={`relative flex items-center justify-center ${className}`}>
      <div className="absolute inset-0 bg-sky-500/20 blur-xl rounded-full" />
      <div className="relative bg-linear-to-tr from-sky-500 to-blue-600 rounded-xl p-2 text-white shadow-lg shadow-sky-500/25 ring-1 ring-white/10">
        <Zap size={size} className="fill-white/20" />
      </div>
    </div>
  );
}

export function LoadingSpinner({ className = "" }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={`animate-spin ${className}`}
    >
      <path d="M21 12a9 9 0 1 1-6.219-8.56" />
    </svg>
  )
}
