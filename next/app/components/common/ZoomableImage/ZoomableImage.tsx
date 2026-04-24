"use client";

import React, { useState, useEffect } from "react";

export interface ZoomableImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
  src: string;
  alt: string;
}

export const ZoomableImage: React.FC<ZoomableImageProps> = ({ src, alt, className, ...props }) => {
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    if (isFullscreen) {
      document.body.style.overflow = "hidden";
      // Small delay to allow display: block to apply before opacity transition
      const timer = setTimeout(() => setIsVisible(true), 10);
      return () => clearTimeout(timer);
    } else {
      setIsVisible(false);
      document.body.style.overflow = "auto";
    }
    return () => {
      document.body.style.overflow = "auto";
    };
  }, [isFullscreen]);

  return (
    <>
      <img
        src={src}
        alt={alt}
        className={`mt-4 rounded-xl border border-gray-200 dark:border-slate-800 shadow-sm cursor-zoom-in hover:shadow-md transition-all duration-200 ${className || ""}`}
        onClick={() => setIsFullscreen(true)}
        {...props}
      />
      
      {isFullscreen && (
        <div
          className={`fixed inset-0 z-[100] flex items-center justify-center p-4 md:p-10 cursor-zoom-out transition-all duration-300 ease-out ${
            isVisible ? "bg-slate-900/90 backdrop-blur-sm opacity-100" : "bg-transparent backdrop-blur-none opacity-0"
          }`}
          onClick={() => {
            setIsVisible(false);
            setTimeout(() => setIsFullscreen(false), 300); // Wait for transition
          }}
        >
          <img
            src={src}
            alt={alt}
            className={`max-w-full max-h-full object-contain rounded-xl shadow-2xl transition-all duration-300 ease-out ${
              isVisible ? "scale-100 translate-y-0" : "scale-95 translate-y-4"
            }`}
             onClick={(e) => e.stopPropagation()} // Prevent clicking image from closing it immediately (optional, but clicking background is better for closing)
          />
          {/* Close button for better UX on mobile */}
          <button 
            className={`absolute top-4 right-4 bg-black/50 hover:bg-black/80 text-white rounded-full p-2 transition-all duration-300 ${
              isVisible ? "opacity-100" : "opacity-0"
            }`}
            onClick={() => {
                setIsVisible(false);
                setTimeout(() => setIsFullscreen(false), 300);
            }}
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
          </button>
        </div>
      )}
    </>
  );
};
