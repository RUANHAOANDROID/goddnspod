"use client"
import { useState, useEffect } from "react";
import type { MetaFunction } from "@remix-run/node";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { Button } from "~/components/ui/button";
import { Label } from "~/components/ui/label";
import { Input } from "~/components/ui/input";
import { defConfig, getCurrentUrl } from "~/config";
import { useToast  } from "~/hooks/use-toast";
export const meta: MetaFunction = () => {
  return [
    { title: "DDNS 配置" },
    { name: "description", content: "DDNS 配置页面" },
  ];
};

export default function Index() {
  const [conf, setConf] = useState(defConfig);
  const [isLoading, setIsLoading] = useState(true);
  const { toast } = useToast()
  const fetchConfig = async () => {
    const currentUrl = getCurrentUrl();
    try {
      setIsLoading(false);
      const response = await fetch(currentUrl+"config");
      if (!response.ok) {
        throw new Error("Network response was not ok");
      }
      const data = await response.json();
      setConf(data);
      console.log("fetch config success",data);
    } catch (error) {
      console.error(error);
    } finally {
      console.log("fetch config finally");
    }
  };

  const submitConfig  = async () => {
    const currentUrl = getCurrentUrl();
    try {
      const response = await fetch(currentUrl+"config", {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(conf)
      });

      if (response.ok) {
        toast({
          variant: "success",
          description: "保存成功",
        })
        console.log('Success:');
      } else {
        console.error('Error:', response.statusText);
        toast({
          variant: "destructive",
          description: "保存失败.",
        })
      }
    } catch (error) {
      console.error('Error:', error);
      toast({
        variant: "destructive",
        description:"保存异常",
      })
    }finally {
        console.log("submit config finally");
    }
  };
  
  useEffect(() => {
    if (isLoading) {
      fetchConfig();
    }
  }, [conf, isLoading]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setConf({ ...conf, [e.target.name]: e.target.value });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log("配置已保存:", conf);
    submitConfig();
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <Card className="w-full max-w-lg">
        <CardHeader className="bg-primary text-primary-foreground">
          <CardTitle>DNSPod DDNS GO</CardTitle>
          <CardDescription>请输入您的DDNS配置信息</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="p-4">
            {Object.entries(conf).map(([key, value]) => (
              <div key={key} className="mb-4">
                <Label htmlFor={key}>{key}</Label>
                <Input
                  id={key}
                  type="text"
                  name={key}
                  value={value}
                  onChange={handleChange}
                />
              </div>
            ))}
          </form>
        </CardContent>
        <CardFooter>
          <Button type="submit" onClick={handleSubmit}>
            保存配置
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
