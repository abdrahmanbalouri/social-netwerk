"use client";
import { createContext, useContext, useEffect, useState } from "react";

const DarkModeContext = createContext(undefined);

export function DarkModeProvider({ children }) {
  const [darkMode, setDarkMode] = useState(false);

  useEffect(() => {
    if (typeof window !== 'undefined') {
      const savedMode = localStorage.getItem("darkMode") === "true";
      setDarkMode(savedMode);
    }
  }, []);

  const toggle = () => {
    setDarkMode(prev => {
      if (typeof window !== 'undefined') {
        localStorage.setItem("darkMode", String(!prev));
      }
      return !prev;
    });
  };

  return (
    <DarkModeContext.Provider value={{ darkMode, toggle }}>
      {children}
    </DarkModeContext.Provider>
  );
}

export function useDarkMode() {
  const context = useContext(DarkModeContext);
  if (context === undefined) {
    throw new Error('useDarkMode must be used within a DarkModeProvider');
  }
  return context;
}