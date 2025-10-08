'use client';

import { useState } from 'react';

interface CodeBlockProps {
  code: string;
  language?: string;
}

export default function CodeBlock({ code, language = 'bash' }: CodeBlockProps) {
  const [copied, setCopied] = useState(false);

  const copyToClipboard = async () => {
    await navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="relative group my-4 sm:my-6">
      <div className="absolute right-2 top-2 z-10">
        <button
          onClick={copyToClipboard}
          className="px-2 sm:px-3 py-1 text-xs bg-gray-700 hover:bg-gray-600 cursor-pointer text-white rounded opacity-100 sm:opacity-0 sm:group-hover:opacity-100 transition-all duration-200"
        >
          {copied ? '✓ Copied!' : 'Copy'}
        </button>
      </div>
      <pre className="bg-gray-900 text-gray-100 p-3 sm:p-4 rounded-lg overflow-x-auto border border-gray-700">
        <code className={`language-${language} text-xs sm:text-sm font-source-code block`}>{code}</code>
      </pre>
    </div>
  );
}

