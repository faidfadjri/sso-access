"use client";


import { useEffect, useState } from "react";
import { TokenExchangeRequest } from "../api/services/oauth/oauth.type";
import { localStorageLib } from "../libs/local-storage";
import { useToast } from "../context/ToastContext";
import { generateAppAuthorizationURL, oauthGetMe, oauthTokenExchange } from "../api/services/oauth/oauth.service";
import { AppCard } from "@/components";
import { authSession } from "../api";
import { getUserAccess } from "../api/services/user-access/user-access.service";
import { UserAccess } from "../api/services/user-access/user-access.service.type";

export default function Dashboard() {

  const { showToast } = useToast();
  const [allApps, setAllApps] = useState<UserAccess[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const handleAuth = async () => {
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get("code");
      const code_verifier = localStorageLib.get<string>("code_verifier");

      if (code) {
        const tokenExchangeRequest: TokenExchangeRequest = {
          code: code,
          redirect_uri: window.location.origin,
          client_id: process.env.NEXT_PUBLIC_CLIENT_ID || "",
          grant_type: "authorization_code",
          code_verifier: code_verifier || "",
        };

        try {
          await oauthTokenExchange(tokenExchangeRequest);
          
          // Remove code from URL to prevent reuse
          window.history.replaceState({}, document.title, window.location.pathname);
          
          const response = await oauthGetMe();
          authSession.save(response.data);

          // Optional: clear code_verifier as it's no longer needed
          localStorageLib.remove("code_verifier");
          window.location.reload();
        } catch (error) {
          console.error(error);
          showToast("An error occurred during authentication", "error");
        }
      } else {
        const currentSession = authSession.get();
        
        if (!currentSession) {
          try {
            const response = await oauthGetMe();
            authSession.save(response.data);
            window.location.reload();
          } catch (error) {
            console.error('Failed to fetch user profile', error);
          }
        } else {
          console.log('Session already exists in local storage');
        }
      }
    };

    handleAuth();

    // Load available apps
    const fetchApps = async () => {
      try {
        setIsLoading(true);
        const response = await getUserAccess({ page: 1, show: -1, user_id: authSession.getUserId() });
        setAllApps(response.data.rows || []);
      } catch (error) {
        console.error("Failed to fetch apps", error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchApps();
  }, [showToast]);

  const handleAppClick = (app: UserAccess) => {
    const storedRecent = localStorageLib.get<UserAccess[]>("recentApps") || [];
    const filtered = storedRecent.filter((a) => a.access_id !== app.access_id);
    const updatedRecent = [app, ...filtered].slice(0, 5); // Keep top 5
    localStorageLib.set("recentApps", updatedRecent);

    const url = generateAppAuthorizationURL(app);
    window.open(url, "_blank", "noopener,noreferrer");
  };

  const getIconPath = (name: string) => {
    if (!name) return "/apps/website.png";
    const lowerName = name.toLowerCase();
    if (lowerName.includes("sivalet")) return "/apps/sivalet.png";
    if (lowerName.includes("promotive")) return "/apps/promotive.png";
    if (lowerName.includes("attendify")) return "/apps/attendify.png";
    if (lowerName.includes("parking")) return "/apps/parking.png";
    return "/apps/website.png";
  };

  return (
    <div className="flex flex-col gap-8">
      <section className="w-full">
        <h2 className="font-sm font-semibold text-gray-500 tracking-wider mb-4">Apps</h2>
        <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 md:gap-6">
          {isLoading ? (
            Array.from({ length: 4 }).map((_, index) => <AppCard.Skeleton key={index} />)
          ) : (
            allApps.map((app, index) => (
              <AppCard
                key={index}
                name={app.service_name}
                iconPath={getIconPath(app.service_name)}
                onClick={() => handleAppClick(app)}
              />
            ))
          )}
        </div>
      </section>

    </div>
  );
}
