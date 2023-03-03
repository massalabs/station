import axios, { AxiosResponse } from 'axios'
import { Plugin } from '../../../shared/interfaces/IPlugin'

export class axiosServices {

  static installPlugin(url: string): Promise<AxiosResponse<number>> {
    return axios.post<any>(`/plugin-manager?source=${url}`)
  }

  static manageLifePlugins(ID: string, command: string): Promise<AxiosResponse<number>> {
    return axios.post<any>(`/plugin-manager/${ID}/execute`, {
      command: command,
    })
  }

  static getPluginsInfo(): Promise<AxiosResponse<Plugin[]>> {
    return axios.get<Plugin[]>(`/plugin-manager`);
  }

  static getpluginInfo(ID: string): Promise<AxiosResponse<{ status: string }>> {
    return axios.get<{ status: string }>(`/plugin-manager/${ID}`);
  }

  static deletePlugins(ID: string): Promise<AxiosResponse<any>> {
    return axios.delete<any>(`/plugin-manager/${ID}`)

  }

  static uploadPlugins(fileName: string): string {
    // const formData = new FormData()
    // formData.append('file', file)
    // return axios.post<any>(`${window.location.hostname}/plugin-manager/upload`, formData)
    return "Not implemented yet"
  }
}

export default axiosServices