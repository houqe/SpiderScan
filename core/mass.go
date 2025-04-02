package core

//var (
//	SrcIP  string           = "10.x.x.x"                                           // 源IP
//	DstIp  string           = "188.131.x.x"                                        // 目标IP
//	device string           = "en0"                                                // 网卡名称
//	SrcMac net.HardwareAddr = net.HardwareAddr{0xf0, 0x18, 0x98, 0x1a, 0x57, 0xe8} // 源mac地址
//	DstMac net.HardwareAddr = net.HardwareAddr{0x5c, 0xc9, 0x99, 0x33, 0x37, 0x80} // 网关mac地址
//)

//
//// 本地状态表的数据结构
//type ScanData struct {
//	ip     string
//	port   int
//	time   int64 // 发送时间
//	retry  int   // 重试次数
//	status int   // 0 未发送 1 已发送 2 已回复 3 已放弃
//}
//
//func recv(datas *[]ScanData, lock *sync.Mutex) {
//	var (
//		snapshot_len int32         = 1024
//		promiscuous  bool          = false
//		timeout      time.Duration = 30 * time.Second
//		handle       *pcap.Handle
//	)
//	handle, _ = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
//	// 使用 handle 作为数据包源来处理所有数据包
//	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
//	scandata := *datas
//
//	for {
//		packet, err := packetSource.NextPacket()
//		if err != nil {
//			continue
//		}
//
//		if IpLayer := packet.Layer(layers.LayerTypeIPv4); IpLayer != nil {
//			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
//				tcp, _ := tcpLayer.(*layers.TCP)
//				ip, _ := IpLayer.(*layers.IPv4)
//				if tcp.Ack != 111223 {
//					continue
//				}
//				if tcp.SYN && tcp.ACK {
//					fmt.Println(ip.SrcIP, " port:", int(tcp.SrcPort))
//					_index := int(tcp.DstPort)
//					lock.Lock()
//					scandata[_index].status = 2
//					lock.Unlock()
//
//				} else if tcp.RST {
//					fmt.Println(ip.SrcIP, " port:", int(tcp.SrcPort), " close")
//					_index := int(tcp.DstPort)
//					lock.Lock()
//					scandata[_index].status = 2
//					lock.Unlock()
//				}
//			}
//		}
//
//		//fmt.Printf("From src port %d to dst port %d\n", tcp.SrcPort, tcp.DstPort)
//	}
//}
//
//func send(index chan int, datas *[]ScanData, lock *sync.Mutex) {
//	srcip := net.ParseIP(SrcIP).To4()
//
//	var (
//		snapshot_len int32 = 1024
//		promiscuous  bool  = false
//		err          error
//		timeout      time.Duration = 30 * time.Second
//		handle       *pcap.Handle
//	)
//	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer handle.Close()
//	scandata := *datas
//	for {
//		_index := <-index
//
//		lock.Lock()
//		data := scandata[_index]
//		port := data.port
//		scandata[_index].status = 1
//		dstip := net.ParseIP(data.ip).To4()
//		lock.Unlock()
//
//		eth := &layers.Ethernet{
//			SrcMAC:       SrcMac,
//			DstMAC:       DstMac,
//			EthernetType: layers.EthernetTypeIPv4,
//		}
//		// Our IPv4 header
//		ip := &layers.IPv4{
//			Version:    4,
//			IHL:        5,
//			TOS:        0,
//			Length:     0, // FIX
//			Id:         0,
//			Flags:      layers.IPv4DontFragment,
//			FragOffset: 0,  //16384,
//			TTL:        64, //64,
//			Protocol:   layers.IPProtocolTCP,
//			Checksum:   0,
//			SrcIP:      srcip,
//			DstIP:      dstip,
//		}
//		// Our TCP header
//		tcp := &layers.TCP{
//			SrcPort:  layers.TCPPort(_index),
//			DstPort:  layers.TCPPort(port),
//			Seq:      111222,
//			Ack:      0,
//			SYN:      true,
//			Window:   1024,
//			Checksum: 0,
//			Urgent:   0,
//		}
//		//tcp.DataOffset = 5 // uint8(unsafe.Sizeof(tcp))
//		_ = tcp.SetNetworkLayerForChecksum(ip)
//		buf := gopacket.NewSerializeBuffer()
//		err := gopacket.SerializeLayers(
//			buf,
//			gopacket.SerializeOptions{
//				ComputeChecksums: true, // automatically compute checksums
//				FixLengths:       true,
//			},
//			eth, ip, tcp,
//		)
//		if err != nil {
//			log.Fatal(err)
//		}
//		//fmt.Println("\n" + hex.EncodeToString(buf.Bytes()))
//		err = handle.WritePacketData(buf.Bytes())
//		if err != nil {
//			fmt.Println(err)
//		}
//	}
//}
//
//func main() {
//	version := pcap.Version()
//	fmt.Println(version)
//	retry := 8
//
//	var datas []ScanData
//	lock := &sync.Mutex{}
//	for i := 20; i < 1000; i++ {
//		temp := ScanData{
//			port:   i,
//			ip:     DstIp,
//			retry:  0,
//			status: 0,
//			time:   time.Now().UnixNano() / 1e6,
//		}
//		datas = append(datas, temp)
//	}
//	fmt.Println("target", DstIp, " count:", len(datas))
//
//	rate := 300
//	distribution := make(chan int, rate)
//
//	go func() {
//		// 每秒将ports数据分配到distribution
//		index := 0
//		for {
//			OldTimestap := time.Now().UnixNano() / 1e6
//
//			for i := index; i < index+rate; i++ {
//				if len(datas) <= index {
//					break
//				}
//				index++
//				distribution <- i
//
//			}
//			if len(datas) <= index {
//				break
//			}
//			Timestap := time.Now().UnixNano() / 1e6
//			TimeTick := Timestap - OldTimestap
//			if TimeTick < 1000 {
//				time.Sleep(time.Duration(1000-TimeTick) * time.Millisecond)
//			}
//		}
//		fmt.Println("发送完毕..")
//	}()
//
//	go recv(&datas, lock)
//	go send(distribution, &datas, lock)
//	// 监控
//	for {
//		time.Sleep(time.Second * 1)
//		count_1 := 0
//		count_2 := 0
//		count_3 := 0
//		var ids []int
//		lock.Lock()
//		for index, data := range datas {
//			if data.status == 1 {
//				count_1++
//				if data.retry >= retry {
//					datas[index].status = 3
//					continue
//				}
//				nowtime := time.Now().UnixNano() / 1e6
//				if nowtime-data.time >= 1000 {
//					datas[index].retry += 1
//					datas[index].time = nowtime
//					ids = append(ids, index)
//					//fmt.Println("重发id:", index)
//					//distribution <- index
//				}
//			} else if data.status == 2 {
//				count_2++
//			} else if data.status == 3 {
//				count_3++
//			}
//		}
//		lock.Unlock()
//		if len(ids) > 0 {
//			time.Sleep(time.Second)
//			increase := 0
//			interval := 60
//			for _, v := range ids {
//				distribution <- v
//				increase++
//				if increase > 1 && increase%interval == 0 {
//					time.Sleep(time.Second)
//				}
//			}
//		}
//		fmt.Println("status=1:", count_1, "status=2:", count_2, "status=3:", count_3)
//	}
//}
