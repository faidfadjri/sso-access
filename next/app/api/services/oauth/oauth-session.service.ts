import { localStorageLib } from "@/app/libs/local-storage";
import type { JWTPayload } from "./oauth.type";

const AUTH_SESSION_KEY = "auth_session";

export const authSession = {
  save(response: JWTPayload) {
    localStorageLib.set<JWTPayload>(AUTH_SESSION_KEY, response);
  },

  get() {
    return localStorageLib.get<JWTPayload>(AUTH_SESSION_KEY);
  },

  getUserId() {
    return this.get()?.user_id;
  },

  getFullName() {
    return this.get()?.full_name;
  },

  getRole() {
    return this.get()?.role_name;
  },

  getService() {
    return this.get()?.service_name;
  },

  getPhoto() {
    return this.get()?.photo;
  },


  clear() {
    localStorageLib.remove(AUTH_SESSION_KEY);
  },

  updateUser(user: Partial<JWTPayload>) {
    const currentSession = this.get();
    if (currentSession) {
      const updatedSession = {
        ...currentSession,
        ...user
      };
      this.save(updatedSession);
    }
  }
};
