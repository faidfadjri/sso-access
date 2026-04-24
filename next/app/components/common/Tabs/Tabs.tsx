"use client";

import React, { useState } from "react";

export interface TabItem {
  label: string;
  content: React.ReactNode;
}

export interface TabsProps {
  items: TabItem[];
  defaultActiveIndex?: number;
}

export const Tabs: React.FC<TabsProps> = ({ items, defaultActiveIndex = 0 }) => {
  const [activeIndex, setActiveIndex] = useState(defaultActiveIndex);

  return (
    <div className="w-full bg-white dark:bg-slate-900 rounded-xl shadow-sm border border-gray-200 dark:border-slate-800 overflow-hidden mt-4 mb-8">
      <div className="flex border-b border-gray-200 dark:border-slate-800 bg-gray-50 dark:bg-slate-800/50">
        {items.map((item, index) => (
          <button
            key={index}
            className={`py-3 px-6 text-sm font-semibold focus:outline-none transition-all duration-200 relative ${
              activeIndex === index
                ? "text-brand-600 dark:text-brand-400 bg-white dark:bg-slate-900"
                : "text-gray-500 hover:text-gray-900 dark:text-slate-400 dark:hover:text-slate-200 hover:bg-gray-100 dark:hover:bg-slate-800"
            }`}
            onClick={() => setActiveIndex(index)}
          >
            {item.label}
            {activeIndex === index && (
              <div className="absolute bottom-0 left-0 w-full h-0.5 bg-brand-500 dark:bg-brand-400" />
            )}
          </button>
        ))}
      </div>
      <div className="p-0 bg-slate-950">
        {items[activeIndex]?.content}
      </div>
    </div>
  );
};
