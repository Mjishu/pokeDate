export const ssr = false;
export const prerender = true;
import { goto } from '$app/navigation';
import type { incomingUser } from '../helpers/users.js';

/**
 * 
 */

export async function load({ fetch }) {
      await GetTokens(fetch);
      let userData: incomingUser | null = await GetCurrentUser(fetch)
      return {
            userData
      };
}

async function GetTokens(fetch: typeof window.fetch): Promise<void> {
      try {
            const refreshToken = localStorage.getItem("refresh_token")
            const bearerToken = "Bearer " + refreshToken
            if (!refreshToken) {
                  console.log("i dont have a refresh token. log in")
                  return
            }
            const fetchParams = {
                  method: "POST",
                  headers: {
                        "Content-Type": "application/json",
                        "Authorization": bearerToken
                  }
            }

            const response = await fetch("/api/refresh", fetchParams)
            const data = await response.json()
            if (data.token) {
                  localStorage.setItem("token", data.token)
            }
      } catch (error) {
            console.error(`error fetching tokens ${error}`)
            return
      }
}

async function GetCurrentUser(fetch: typeof window.fetch): Promise<incomingUser | null> {
      try {
            const token = localStorage.getItem("token")
            const bearerToken = "Bearer " + token
            const fetchParams = {
                  method: "POST",
                  headers: {
                        "Content-Type": "application/json",
                        "Authorization": bearerToken
                  }
            }
            const response = await fetch("/api/users/current", fetchParams)
            const data = await response.json()
            if (response.status == 200) {
                  console.log(data)
                  return data
            } else {
                  return null
            }
      } catch (err) {
            console.error(`error fetching curernt user data ${err}`)
            return null
      }
}