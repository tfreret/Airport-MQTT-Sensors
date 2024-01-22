import {
  DateRangePicker,
  DateRangePickerValue,
  Flex,
  Icon,
  LineChart,
  Title
} from "@tremor/react";
import {useState, useEffect} from "react";
import axios from 'axios';

interface Point {
  date: string
  value: number
}

interface PointAPI {
  Time: string,
  Value: number
}

import { InformationCircleIcon } from "@heroicons/react/solid";

interface Props {
  url: string
}

const initialFrom = new Date(0);

export default function Chart({url}: Props) {

  const [chartData, setChartData] = useState<Point[]>([]);
  const [selectedDate, setSelectedDate] = useState<DateRangePickerValue>({from: undefined, to: undefined});
  const fetchData = async () => {
    try {
      const from = selectedDate.from ? selectedDate.from.toISOString() : initialFrom.toISOString();
      const to = selectedDate.to ? selectedDate.to.toISOString() : new Date().toISOString();
      const response = await axios.get(`${url}?from=${from}&to=${to}`);
        const data = response.data['Points'] as PointAPI[];
        const points = data === null ? [] : data.map(it=> {return {date: it.Time, value: it.Value}});
        setChartData(
          points
        );

    } catch (error) {
      setChartData([]);
      console.error('Erreur lors de la récupération des données:', error);
    }
  };

  useEffect(() => {
    fetchData();
    const intervalId = setInterval(() => {
        fetchData()
      }, 1000 * 5) // in milliseconds
      return () => clearInterval(intervalId)
  }, [selectedDate]);

  return (
    <>
      <div className="md:flex justify-between">
        <div>
          <Flex className="space-x-0.5" justifyContent="start" alignItems="center">
            <Title> Données du capteur </Title>
            <Icon
              icon={InformationCircleIcon}
              variant="simple"
              tooltip="Afficher les données du capteur"
            />
          </Flex>
        </div>
        <div>
          <DateRangePicker value={selectedDate} onValueChange={setSelectedDate}/>
        </div>
      </div>
      <div className="mt-8 hidden sm:block">
        <LineChart
            className="mt-5 h-72"
            data={chartData}
            index="date"
            categories={['value']}
            showLegend={false}
            yAxisWidth={60}
            showXAxis
            showGridLines
            autoMinValue
            colors={['blue', 'indigo', 'rose-400', '#ffcc33']}
    />
      </div>
    </>
  );
}