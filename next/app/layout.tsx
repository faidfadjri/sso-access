import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { ToastProvider } from "@/app/context/ToastContext";
import { SidebarProvider } from "./components/layouts";
import NextTopLoader from 'nextjs-toploader';

const inter = Inter({
  variable: "--font-inter",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Akastra Access",
  description: "Identity Provider for Akastra Toyota",
};

import { TopbarProvider } from "./context/TopbarContext";

import { AppShell } from "./components/layouts";

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head><link rel="icon" href="/favicon.png" type="image/png" sizes="any" /></head>
      <body
        className={`${inter.variable} antialiased`}
        suppressHydrationWarning
      >
        <NextTopLoader color="#0975A1" showSpinner={false} />
        <TopbarProvider>
          <SidebarProvider>
            <ToastProvider>
              <AppShell>
                {children}
              </AppShell>
            </ToastProvider>
          </SidebarProvider>
        </TopbarProvider>
      </body>
    </html>
  );
}
