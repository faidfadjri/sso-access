"use client";

import React, { createContext, useContext, useState, ReactNode, useEffect } from "react";

interface TopbarContextType {
  title: ReactNode;
  setTitle: (title: ReactNode) => void;
}

const TopbarContext = createContext<TopbarContextType | undefined>(undefined);

export function TopbarProvider({ children }: { children: ReactNode }) {
  const [title, setTitle] = useState<ReactNode>(null);

  return (
    <TopbarContext.Provider value={{ title, setTitle }}>
      {children}
    </TopbarContext.Provider>
  );
}

export function useTopbar() {
  const context = useContext(TopbarContext);
  if (context === undefined) {
    throw new Error("useTopbar must be used within a TopbarProvider");
  }
  return context;
}

export function TopbarSetter({ title }: { title: ReactNode }) {
  const { setTitle } = useTopbar();

  useEffect(() => {
    setTitle(title);
    // Cleanup title on unmount if necessary, or let next page overwrite it
  }, [title, setTitle]);

  return null;
}
