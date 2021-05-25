import { useState, useEffect } from "react";
import { subMinutes } from "date-fns";
import axios from "axios";

export type Data = {
  Ip: string;
  RequestId: string;
  type: string;
  UserId: number;
  category: "";
  debugId: "";
  id: number;
  level: "error" | "info" | "warning";
  msg: string;
  time: string;
}[];

export const useData = () => {
  const ENDPOINT = "http://:3000/logs";
  //const ENDPOINT = "http://91.192.221.250:3000/logs";
  const [data, setData] = useState<Data>([]);
  const [page, setPage] = useState(1);

  const last = page * 10;
  const first = last - 10;

  useEffect(() => {
    (async () => {
      try {
        const res = await axios.get(ENDPOINT);
        setData(res.data);
      } catch (error) {
        console.log(error);
      }
    })();
  }, []);

  return {
    data: data.slice(first, last),
    setPage,
    totalPages: Math.round(data.length / 10),
  };
};
