"use client";

import { Inter, Source_Code_Pro } from "next/font/google";
import "./globals.css";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import { useState, useCallback } from "react";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
});

const source_code_pro = Source_Code_Pro({
  subsets: ["latin"],
  variable: "--font-source-code",
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const toggleSidebar = useCallback(() => {
    setIsSidebarOpen((prev) => !prev);
  }, []);

  const closeSidebar = useCallback(() => {
    setIsSidebarOpen(false);
  }, []);

  return (
    <html lang="en" className={`dark ${inter.variable} ${source_code_pro.variable}`}>
      <head>
        <title>starknode-kit Documentation</title>
        <meta
          name="description"
          content="Complete documentation for starknode-kit - A CLI tool for setting up and managing Ethereum and Starknet nodes"
        />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
      </head>
      <body className="antialiased bg-gray-900 text-gray-100">
        <Sidebar isOpen={isSidebarOpen} onClose={closeSidebar} />
        <Header onMenuClick={toggleSidebar} />
        <main className="lg:ml-64 pt-16 min-h-screen bg-gray-900">
          <div className="max-w-5xl mx-auto px-4 sm:px-6 md:px-8 lg:px-12 py-8 sm:py-12 lg:py-16">
            {children}
          </div>
        </main>
      </body>
    </html>
  );
}
