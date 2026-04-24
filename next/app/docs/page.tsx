import React from "react";
import { Tabs, ZoomableImage } from "@/app/components/common";
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';
import Link from "next/link";

export const metadata = {
  title: "Documentation | Akastra Access",
  description: "Learn how to integrate with Akastra Access Identity Provider.",
};

export default function DocsPage() {
  const redirectJsCode = `// Example: Redirecting user with PKCE (Proof Key for Code Exchange)
async function authorize() {
  // 1. Generate state and code verifier
  const state = Math.random().toString(36).substring(2, 15);
  const codeVerifier = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15);
  
  // 2. Generate code challenge (SHA-256)
  const data = new TextEncoder().encode(codeVerifier);
  const digest = await crypto.subtle.digest("SHA-256", data);
  const codeChallenge = btoa(String.fromCharCode(...new Uint8Array(digest)))
    .replace(/\\+/g, "-").replace(/\\//g, "_").replace(/=+$/, "");

  // 3. Save verifier to localStorage to use later in callback
  localStorage.setItem("code_verifier", codeVerifier);

  // 4. Build authorization URL
  const queryParams = new URLSearchParams({
    client_id: "YOUR_CLIENT_ID",
    redirect_uri: "http://localhost:3000/callback",
    response_type: "code",
    scope: "all",
    state: state,
    code_challenge: codeChallenge,
    code_challenge_method: "S256",
  });

  // 5. Redirect user
  window.location.href = "http://your-identity-provider.com/api/v1/oauth/authorize?" + queryParams.toString();
}`;

  const redirectGoCode = `package main

import (
\t"crypto/rand"
\t"encoding/hex"
\t"net/http"
\t"net/url"

\t"github.com/go-chi/chi/v5"
)

func generateState() string {
\tbytes := make([]byte, 16)
\trand.Read(bytes)
\treturn hex.EncodeToString(bytes)
}

// Example: Redirecting user with go-chi
func loginHandler(w http.ResponseWriter, r *http.Request) {
\t// 1. Generate state
\tstate := generateState()

\t// 2. Save state to session/cookie to verify later in callback
\t// http.SetCookie(w, &http.Cookie{Name: "oauth_state", Value: state})

\t// 3. Build authorization URL
\tproviderURL := "http://your-identity-provider.com/api/v1/oauth/authorize"
\tu, _ := url.Parse(providerURL)
\tq := u.Query()
\tq.Set("client_id", "YOUR_CLIENT_ID")
\tq.Set("redirect_uri", "http://localhost:3000/callback")
\tq.Set("response_type", "code")
\tq.Set("scope", "all")
\tq.Set("state", state)
\tu.RawQuery = q.Encode()

\t// 4. Redirect the user
\thttp.Redirect(w, r, u.String(), http.StatusFound)
}

func main() {
\tr := chi.NewRouter()
\tr.Get("/login", loginHandler)
\thttp.ListenAndServe(":3000", r)
}`;

  const callbackJsCode = `import { useEffect } from "react";

// Example: Handling callback in React
export function useOAuthCallback() {
  useEffect(() => {
    const handleAuth = async () => {
      // 1. Get code from URL and verifier from storage
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get("code");
      const code_verifier = localStorage.getItem("code_verifier");

      if (!code || !code_verifier) return;

      try {
        // 2. Perform token exchange (PKCE)
        const response = await fetch("http://your-identity-provider.com/api/v1/oauth/token", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            client_id: "YOUR_CLIENT_ID",
            grant_type: "authorization_code",
            redirect_uri: window.location.origin,
            code: code,
            code_verifier: code_verifier,
          }),
        });

        const data = await response.json();
        console.log("Access Token:", data.access_token);

        // 3. Clean up URL and storage to prevent reuse
        window.history.replaceState({}, document.title, window.location.pathname);
        localStorage.removeItem("code_verifier");

        // 4. Reload or redirect to dashboard
        window.location.reload();
      } catch (error) {
        console.error("Token exchange failed:", error);
      }
    };

    handleAuth();
  }, []);
}`;

  const callbackGoCode = `package main

import (
\t"bytes"
\t"encoding/json"
\t"net/http"
)

// Example: Handling callback in go-chi
func callbackHandler(w http.ResponseWriter, r *http.Request) {
\t// 1. Get code from URL
\tcode := r.URL.Query().Get("code")
\tif code == "" {
\t\thttp.Error(w, "Code not found", http.StatusBadRequest)
\t\treturn
\t}

\t// 2. Perform token exchange (using client_secret)
\ttokenURL := "http://your-identity-provider.com/api/v1/oauth/token"
\tpayload := map[string]interface{}{
\t\t"client_id":     "YOUR_CLIENT_ID",
\t\t"client_secret": "YOUR_CLIENT_SECRET",
\t\t"grant_type":    "authorization_code",
\t\t"redirect_uri":  "http://localhost:3000/callback",
\t\t"code":          code,
\t}
\tbody, _ := json.Marshal(payload)

\tresp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(body))
\tif err != nil || resp.StatusCode != http.StatusOK {
\t\thttp.Error(w, "Failed to exchange token", http.StatusInternalServerError)
\t\treturn
\t}
\tdefer resp.Body.Close()

\tvar data map[string]interface{}
\tjson.NewDecoder(resp.Body).Decode(&data)

\t// 3. Optional: Verify state from cookie/session here

\t// 4. Redirect to dashboard or use token
\taccessToken, _ := data["access_token"].(string)
\tw.Write([]byte("Access Token: " + accessToken))
}`;

  const redirectPhpCode = `<?php

namespace App\\Http\\Controllers;

use Illuminate\\Http\\Request;
use Illuminate\\Support\\Str;

class AuthController extends Controller
{
    // Example: Redirecting user in Laravel
    public function redirect()
    {
        // 1. Generate state
        $state = Str::random(40);

        // 2. Save state to session to verify later in callback
        session(['oauth_state' => $state]);

        // 3. Build authorization URL
        $query = http_build_query([
            'client_id' => 'YOUR_CLIENT_ID',
            'redirect_uri' => 'http://localhost:3000/callback',
            'response_type' => 'code',
            'scope' => 'all',
            'state' => $state,
        ]);

        $providerUrl = 'http://your-identity-provider.com/api/v1/oauth/authorize';

        // 4. Redirect the user
        return redirect()->away($providerUrl . '?' . $query);
    }
}`;

  const callbackPhpCode = `<?php

namespace App\\Http\\Controllers;

use Illuminate\\Http\\Request;
use Illuminate\\Support\\Facades\\Http;

class AuthController extends Controller
{
    // ... Redirect method ...

    // Example: Handling callback in Laravel
    public function callback(Request $request)
    {
        // 1. Get code from URL
        $code = $request->query('code');
        if (!$code) {
            return response('Code not found', 400);
        }

        // 2. Perform token exchange (using client_secret)
        $response = Http::post('http://your-identity-provider.com/api/v1/oauth/token', [
            'client_id' => 'YOUR_CLIENT_ID',
            'client_secret' => 'YOUR_CLIENT_SECRET',
            'grant_type' => 'authorization_code',
            'redirect_uri' => 'http://localhost:3000/callback',
            'code' => $code,
        ]);

        if ($response->failed()) {
            return response('Failed to exchange token', 500);
        }

        $data = $response->json();

        // 3. Optional: Verify state from session here

        // 4. Redirect or save token
        return response('Access Token: ' . $data['access_token']);
    }
}`;

  const CodeBlock = ({ code, language }: { code: string, language: string }) => (
    <div className="rounded-xl overflow-hidden bg-[#1e1e1e] shadow-2xl border border-gray-700/50 mt-2 mb-4">
      <div className="flex items-center px-4 py-3 bg-[#2d2d2d] border-b border-gray-700/50 relative">
        <div className="flex space-x-2 absolute">
          <div className="w-3 h-3 rounded-full bg-[#ff5f56]" />
          <div className="w-3 h-3 rounded-full bg-[#ffbd2e]" />
          <div className="w-3 h-3 rounded-full bg-[#27c93f]" />
        </div>
        <div className="w-full text-center text-xs text-gray-400 font-mono select-none">
          {language === "javascript" ? "script.js" : language === "php" ? "index.php" : "main.go"}
        </div>
      </div>
      <div className="p-4 overflow-x-auto text-sm font-mono">
        <SyntaxHighlighter
          language={language}
          style={vscDarkPlus}
          customStyle={{ background: 'transparent', padding: 0, margin: 0 }}
          wrapLongLines={false}
        >
          {code}
        </SyntaxHighlighter>
      </div>
    </div>
  );

  return (
    <div className="w-full mx-auto p-5 font-sans">
      <div className="mb-10">
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-3">
          Developer Documentation
        </h1>
        <p className="text-sm opacity-70 italic">© 2026 <Link href="https://faidfadjri.space" className="underline font-semibold" target="_blank">Faid Fadjri</Link> All rights reserved.</p>
      </div>

      <div className="bg-white dark:bg-slate-900 rounded-xl shadow-sm border border-gray-200 dark:border-slate-800 p-8 mb-8">
        <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">
          Step 1: Create Credentials & Identity
        </h2>
        <div className="text-gray-600 dark:text-gray-300 space-y-4">
          <p>
            Before you can authenticate users, you must register your application to get the
            client credentials. You need login as <strong>Super Admin</strong> to create credentials.
          </p>
          <ul className="list-disc list-inside space-y-2 ml-2">
            <li>
              Navigate to the <strong>Clients</strong> menu in the Identity Provider dashboard.
              <br />
              <ZoomableImage src="/docs/client-menu.png" alt="client-menu" className="max-w-md w-full" />
            </li>
            <li className="mt-4">
              Click on <strong>Create New Credentials</strong>.
              <br />
              <ZoomableImage src="/docs/credential-modal.png" alt="credential-modal" className="max-w-md w-full" />
            </li>
            <li>
              Fill in your application name and configure the required details. Once created, you will be provided with a <code>client_id</code> and a <code>client_secret</code>.
            </li>
            <li>
              Safeguard your <code>client_secret</code> safely as it should never be exposed in public-facing applications (like pure SPAs without a backend).
            </li>
          </ul>
        </div>
      </div>

      <div className="bg-white dark:bg-slate-900 rounded-xl shadow-sm border border-gray-200 dark:border-slate-800 p-8 mb-8">
        <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">
          Step 2: Define Redirect URI
        </h2>
        <div className="text-gray-600 dark:text-gray-300 space-y-4">
          <p>
            The Redirect URI is the location where the Identity Provider will send the user back to
            after an authentication attempt. It is critical for security that this matches your application actual callback URL.
          </p>
          <ul className="list-disc list-inside space-y-2 ml-2">
            <li>
              Inside your Client configuration dashboard, locate the <strong>Redirect URIs</strong> section.
            </li>
            <li>
              Add the exact URL where your application handles the OAauth callback (e.g., <code>http://localhost:3000/callback</code>).
            </li>
          </ul>
        </div>
      </div>

      <div className="bg-white dark:bg-slate-900 rounded-xl shadow-sm border border-gray-200 dark:border-slate-800 p-8 mb-8">
        <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">
          Step 3: Frontend Redirect (Authorization Request)
        </h2>
        <div className="text-gray-600 dark:text-gray-300 space-y-4 mb-6">
          <p>
            To begin the authentication flow, you must redirect your user to the Identity Provider's{" "}
            <code>/authorize</code> endpoint along with your <code>client_id</code>, <code>redirect_uri</code>, and <code>response_type=code</code>.
          </p>
          
          <div className="flex flex-col mb-4">
            <span className="text-sm font-semibold text-gray-900 dark:text-white mb-2">Endpoint</span>
            <code className="px-3 py-2 bg-gray-100 dark:bg-slate-800 border border-gray-200 dark:border-slate-700 rounded-md text-sm font-mono text-gray-800 dark:text-gray-200 w-fit">
              GET /api/v1/oauth/authorize
            </code>
          </div>

          <div className="overflow-x-auto my-6 border border-gray-200 dark:border-slate-800 rounded-lg">
            <table className="min-w-full divide-y divide-gray-200 dark:divide-slate-800 text-sm text-left">
              <thead className="bg-gray-50 dark:bg-slate-800/50 text-gray-700 dark:text-gray-300 font-semibold">
                <tr>
                  <th scope="col" className="px-6 py-3">Parameter</th>
                  <th scope="col" className="px-6 py-3">Type</th>
                  <th scope="col" className="px-6 py-3">Required</th>
                  <th scope="col" className="px-6 py-3">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 dark:divide-slate-800 bg-white dark:bg-slate-900 text-gray-600 dark:text-gray-400">
                <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">client_id</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4"><span className="px-2 py-1 text-xs font-medium text-green-700 bg-green-100 rounded-md dark:bg-green-900/30 dark:text-green-400">Yes</span></td>
                  <td className="px-6 py-4">The client ID obtained from the credentials dashboard.</td>
                </tr>
                <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">redirect_uri</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4"><span className="px-2 py-1 text-xs font-medium text-green-700 bg-green-100 rounded-md dark:bg-green-900/30 dark:text-green-400">Yes</span></td>
                  <td className="px-6 py-4">The registered callback URL for your application.</td>
                </tr>
                <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">response_type</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4"><span className="px-2 py-1 text-xs font-medium text-green-700 bg-green-100 rounded-md dark:bg-green-900/30 dark:text-green-400">Yes</span></td>
                  <td className="px-6 py-4">Must be set to <code className="text-gray-800 dark:text-gray-200 bg-gray-100 dark:bg-slate-800 px-1 py-0.5 rounded">code</code>.</td>
                </tr>
                  <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">scope</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4"><span className="px-2 py-1 text-xs font-medium text-green-700 bg-green-100 rounded-md dark:bg-green-900/30 dark:text-green-400">Yes</span></td>
                  <td className="px-6 py-4">Must be set to <code className="text-gray-800 dark:text-gray-200 bg-gray-100 dark:bg-slate-800 px-1 py-0.5 rounded">read+write</code>.</td>
                </tr>
                  <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">state</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4"><span className="px-2 py-1 text-xs font-medium text-green-700 bg-green-100 rounded-md dark:bg-green-900/30 dark:text-green-400">Yes</span></td>
                  <td className="px-6 py-4">Random string.</td>
                </tr>
                <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">code_challenge</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4">
                    <span className="px-2 py-1 text-xs font-medium text-purple-700 bg-purple-100 rounded-md dark:bg-green-900/30 dark:text-green-400">
                     Yes, if PKCE / Frontend
                    </span>
                  </td>
                  <td className="px-6 py-4">Generated random string.</td>
                </tr>
                <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">code_challenge_method</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4">
                    <span className="px-2 py-1 text-xs font-medium text-purple-700 bg-purple-100 rounded-md dark:bg-green-900/30 dark:text-green-400  ">
                     Yes, if PKCE / Frontend
                    </span>
                  </td>
                  <td className="px-6 py-4">Must be set to <code className="text-gray-800 dark:text-gray-200 bg-gray-100 dark:bg-slate-800 px-1 py-0.5 rounded">S256</code>.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <Tabs
          items={[
            {
              label: "JavaScript (Frontend)",
              content: <CodeBlock code={redirectJsCode} language="javascript" />,
            },
            {
              label: "Go (Go-Chi)",
              content: <CodeBlock code={redirectGoCode} language="go" />,
            },
            {
              label: "PHP (Laravel)",
              content: <CodeBlock code={redirectPhpCode} language="php" />,
            },
          ]}
        />
      </div>

      <div className="bg-white dark:bg-slate-900 rounded-xl shadow-sm border border-gray-200 dark:border-slate-800 p-8 mb-8">
        <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">
          Step 4: Callback & Token Exchange
        </h2>
        <div className="text-gray-600 dark:text-gray-300 space-y-4 mb-6">
          <p>
            When the user successfully authenticates, they are redirected back to your{" "}
            <code>redirect_uri</code> with a <code>code</code> query parameter. You must capture this code and exchange it for an <strong>access token</strong> by making a POST request to the Identity Provider's token endpoint.
          </p>
          
          <div className="flex flex-col mb-4">
            <span className="text-sm font-semibold text-gray-900 dark:text-white mb-2">Endpoint</span>
            <code className="px-3 py-2 bg-gray-100 dark:bg-slate-800 border border-gray-200 dark:border-slate-700 rounded-md text-sm font-mono text-gray-800 dark:text-gray-200 w-fit">POST /api/v1/oauth/token</code>
          </div>
          
          <p className="p-4 bg-blue-50 dark:bg-sky-900/30 text-blue-800 dark:text-sky-300 rounded-lg text-sm">
            <strong>Security Note:</strong> The token exchange requires your <code>client_secret</code>. This request should always occur on your secure backend server, not directly from the browser, to avoid leaking your secret.
          </p>

          <div className="overflow-x-auto my-6 border border-gray-200 dark:border-slate-800 rounded-lg">
            <table className="min-w-full divide-y divide-gray-200 dark:divide-slate-800 text-sm text-left">
              <thead className="bg-gray-50 dark:bg-slate-800/50 text-gray-700 dark:text-gray-300 font-semibold">
                <tr>
                  <th scope="col" className="px-6 py-3">Parameter / Payload</th>
                  <th scope="col" className="px-6 py-3">Type</th>
                  <th scope="col" className="px-6 py-3">Required</th>
                  <th scope="col" className="px-6 py-3">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-200 dark:divide-slate-800 bg-white dark:bg-slate-900 text-gray-600 dark:text-gray-400">
                <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">grant_type</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4"><span className="px-2 py-1 text-xs font-medium text-green-700 bg-green-100 rounded-md dark:bg-green-900/30 dark:text-green-400">Yes</span></td>
                  <td className="px-6 py-4">Must be set to <code className="text-gray-800 dark:text-gray-200 bg-gray-100 dark:bg-slate-800 px-1 py-0.5 rounded">authorization_code</code>.</td>
                </tr>
                <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">code</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4"><span className="px-2 py-1 text-xs font-medium text-green-700 bg-green-100 rounded-md dark:bg-green-900/30 dark:text-green-400">Yes</span></td>
                  <td className="px-6 py-4">The authorization code received from the query parameter.</td>
                </tr>
                <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">client_id</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4"><span className="px-2 py-1 text-xs font-medium text-green-700 bg-green-100 rounded-md dark:bg-green-900/30 dark:text-green-400">Yes</span></td>
                  <td className="px-6 py-4">The client ID obtained from the credentials dashboard.</td>
                </tr>
                <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">client_secret</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4">
                    <span className="px-2 py-1 text-xs font-medium text-purple-700 bg-purple-100 rounded-md dark:bg-purple-900/30 dark:text-purple-400 whitespace-nowrap">
                      Yes, if Backend Request
                    </span>
                  </td>
                  <td className="px-6 py-4">The client secret. Required if request comes from backend level.</td>
                </tr>
                <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">redirect_uri</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4"><span className="px-2 py-1 text-xs font-medium text-green-700 bg-green-100 rounded-md dark:bg-green-900/30 dark:text-green-400">Yes</span></td>
                  <td className="px-6 py-4">The exact callback URL that was used in the authorization request.</td>
                </tr>
                <tr>
                  <td className="px-6 py-4 font-mono text-xs text-blue-600 dark:text-blue-400">code_verifier</td>
                  <td className="px-6 py-4">string</td>
                  <td className="px-6 py-4">
                    <span className="px-2 py-1 text-xs font-medium text-purple-700 bg-purple-100 rounded-md dark:bg-purple-900/30 dark:text-purple-400 whitespace-nowrap">
                      Yes, if PKCE / Frontend
                    </span>
                  </td>
                  <td className="px-6 py-4">The original unhashed verifier string. Required if using PKCE from the frontend.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <Tabs
          items={[
            {
              label: "JavaScript (Node/Axios)",
              content: <CodeBlock code={callbackJsCode} language="javascript" />,
            },
            {
              label: "Go (Go-Chi)",
              content: <CodeBlock code={callbackGoCode} language="go" />,
            },
            {
              label: "PHP (Laravel)",
              content: <CodeBlock code={callbackPhpCode} language="php" />,
            },
          ]}
        />
      </div>
    </div>
  );
}
