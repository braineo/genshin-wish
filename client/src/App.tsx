import React, { useEffect, useState, useRef } from "react";
import { useAxios } from "./axios";
import styles from "./App.module.less";
import { Tabs } from "antd";
import { Route, useHistory, useParams } from "react-router";
const { TabPane } = Tabs;

type GachaConfig = {
  name: string;
  key: string;
};

function App() {
  const client = useAxios();
  const [gachaConfigs, setGachaConfigs] = useState<GachaConfig[]>([]);
  const params = useParams<{ configKey: string }>();
  const history = useHistory();
  useEffect(() => {
    const fetchConfigs = async () => {
      const gachaLog = await client.get<{ data: GachaConfig[] }>("gacha");
      setGachaConfigs(gachaLog.data.data);
    };
    fetchConfigs();
  }, []);

  const handleTabChange = (key: string) => {
    history.push(key);
  };

  console.log(params.configKey);
  const tabs = [{ name: "全部", key: "all" }, ...gachaConfigs];

  return (
    <Tabs onChange={handleTabChange} type="card" activeKey={params.configKey}>
      {tabs.map((config) => (
        <TabPane tab={config.name} key={config.key}>
          Content of Tab Pane 1
        </TabPane>
      ))}
    </Tabs>
  );
}

export default App;
