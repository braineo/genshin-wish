import { Col, List, Form, Input, Button } from 'antd';
import React, { useEffect, useState } from 'react';
import { useHistory } from 'react-router';
import { useAxios } from '../../utils/axios';

type User = {
  id: string;
  name: string;
};
const Home: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);
  const client = useAxios();
  const history = useHistory();
  const [form] = Form.useForm();

  const [loading, setLoading] = useState(false);
  useEffect(() => {
    const fetchUser = async () => {
      try {
        const userData = await client.get<{ data: User[] }>('/user');
        setUsers(userData.data.data);
      } catch (error) {
        console.log(error);
      }
    };
    fetchUser();
  }, []);

  const handleSubmit = async (values: unknown) => {
    setLoading(true);
    try {
      await client.post('/log', values);
    } catch (error) {
    } finally {
      setLoading(false);
    }
  };

  return (
    <Col span={12} offset={6}>
      <Form form={form} name="wish-url" onFinish={handleSubmit}>
        <Form.Item
          name="query"
          label="抽卡记录URL"
          rules={[{ required: true }]}
        >
          <Input name="query" />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit" loading={loading}>
            查询
          </Button>
        </Form.Item>
      </Form>
      <List
        dataSource={users}
        renderItem={item => (
          <List.Item
            actions={[
              <Button onClick={() => history.push(`/stat/${item.id}/all`)}>
                查看
              </Button>,
            ]}
          >
            <List.Item.Meta title={item.name} description={item.id} />
          </List.Item>
        )}
      />
    </Col>
  );
};

export default Home;
