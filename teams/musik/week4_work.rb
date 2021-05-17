require 'rubygems'
require 'grpc'
#require flow-ruby (https://github.com/cybercent/flow-ruby)
require 'flow/access/access_services_pb'
require 'flow/execution/execution_services_pb'
require 'json'

class FlowClient
  def initialize(node_address)
    @stub = Access::AccessAPI::Stub.new(node_address, :this_channel_is_insecure)
  end

  def ping
    req = Access::PingRequest.new
    res = @stub.ping(req)
    res
  end

  def request method,args={}
    req = Object.const_get(["Access::",method.to_s.split('_').collect(&:capitalize).join,"Request"].join).new args
    @stub.send method,req
  end

  def get_account_at_latest_block(address)
    res = request :get_latest_block,{address: to_bytes(address)}
    res.account 
  end

  def get_latest_block is_sealed: true
    res = request :get_latest_block,{is_sealed: is_sealed}
    res.block
  end  

  def execute_script(script, args = [])
    req = Access::ExecuteScriptAtLatestBlockRequest.new(script: script, arguments: args)
    res = @stub.execute_script_at_latest_block(req)
    parse_json(res.value)
  end

  # private

  def parse_json(event_payload)
    JSON.parse(event_payload, object_class: OpenStruct)
  end

  def to_bytes(string)
    [string].pack('H*')
  end

  def to_string(bytes)
    bytes.unpack('H*').first
  end
end
class NbaShot
  def initialize
    @client = FlowClient.new 'access.mainnet.nodes.onflow.org:9000'
    @id = "A.c1e4f4f4c4257510"
  end
  def get_latest_moment_purchased blocks_num = 30
    block = @client.get_latest_block 
    puts "Current Block Height: #{block.height}"
    events = @client.request :get_events_for_height_range,{
        type: "#{@id}.Market.MomentPurchased",
        start_height: block.height-blocks_num,
        end_height: block.height
    }
    txs = events.results.collect{|result|
      result.events.collect{|event|
        parse_event(event)
      }
    }.flatten
    puts "Found txs: #{txs.count}"
    txs
  end
  def parse_event event
     fields = @client.parse_json(event.payload).value.fields
     item = OpenStruct.new fields.collect{|field|
      [field.name,field.value.value]
     }.to_h
     item.id = item.id.to_i
     item.price = item.price.to_f
     item.seller = item.seller.value.to_s
     item.tx_id = @client.to_string(event.transaction_id)
     item.tx_index = event.transaction_index
     item.to_h
  end
end
p NbaShot.new.get_latest_moment_purchased()
#Example outputs
#Current Block Height: 14538215
#Found txs: 32
#[{:id=>5434811, :price=>41.0, :seller=>"0x4de1db695efa1966", :tx_id=>"b55218c43f4dfdfc395041c24293125e44a230cc641f7564011b738befbb2842", :tx_index=>1}, {:id=>9295610, :price=>4.0, :seller=>"0x7693c3515e9bc4fd", :tx_id=>"0d130a3b57a541128e51bd0d2e3827b4aac227bb59da139855c57d4058931905", :tx_index=>1}, {:id=>11124725, :price=>4.0, :seller=>"0x2dc5954366200389", :tx_id=>"73d76f2512d4ef7871dd5a931751a03262055825fc5d0cca9e9314a84971ec4d", :tx_index=>4}, {:id=>8561381, :price=>4.0, :seller=>"0x8aece936034712b3", :tx_id=>"830dfd2f84fb617547407a3b3eb10185756730aa784a80177ef560ff8d0ff5dd", :tx_index=>1}, {:id=>4360109, :price=>5.0, :seller=>"0xd8dddd4d9b5f12e2", :tx_id=>"86f400a90fa78ebf6e26b7168b61b6311eb631ca4752f7de0483d967d0e2ac04", :tx_index=>4}, {:id=>899887, :price=>95.0, :seller=>"0x3d2166ceb390990f", :tx_id=>"a3a3ec7e4838cdfd22ea268a21e61ebc2b1861aacd5efeac1c5eb3c02c70c5ad", :tx_index=>3}, {:id=>9804030, :price=>8.0, :seller=>"0x282fc53e95c854af", :tx_id=>"6ef315c4ce9f5ca4f8dcf4fc76a8d5b88e3da06608d095d3eec73540ca004b6b", :tx_index=>1}, {:id=>6099140, :price=>10.0, :seller=>"0xb06cfd780351a4bd", :tx_id=>"b8d98b5cd8fdc78decd3bdb4882de8ea812bc501f154beb054886c805f2e8b87", :tx_index=>1}, {:id=>9828084, :price=>6.0, :seller=>"0xce169fdaed4dbd84", :tx_id=>"3e7c50db232bc72f5b8406e8f5fc7e842e8bb5a1158732475a8d2aa78bcc2cac", :tx_index=>10}, {:id=>2098177, :price=>34.0, :seller=>"0x86130bbde1981121", :tx_id=>"4099122f77e93ab3b660f11eb99a64a986d8e7a85e432c0d9e886f9dd4947d29", :tx_index=>4}, {:id=>3348635, :price=>8.0, :seller=>"0xb64cfb37bb8a73bf", :tx_id=>"42c219d518b8f15e3de9befb93358ed333a31cb2b20e35ea49991f503877672b", :tx_index=>5}, {:id=>6818260, :price=>4.0, :seller=>"0xc9360e8ebdfadc2f", :tx_id=>"74d3829e4ad6253faf8c5ae790e3be00a212538b6d149c99ee291ed3cb081c65", :tx_index=>6}, {:id=>5014748, :price=>3.0, :seller=>"0x562c73b602f8f5e4", :tx_id=>"901d471060bec2a523e828be0ad431f88efd63756a946ba4e2be95e239252817", :tx_index=>9}, {:id=>10068668, :price=>6.0, :seller=>"0x5000e6e834163a08", :tx_id=>"c12535128980705481eb34f3f180ec17e3b584655ae7c7b2d037b951fa71e9a0", :tx_index=>2}, {:id=>9909827, :price=>4.0, :seller=>"0xd76bd22e0e97d1b5", :tx_id=>"78f1265f79e63a4c68121da19f3d78a700dce5c2647e3f4b0c921aea0d8e95d9", :tx_index=>0}, {:id=>7572341, :price=>4.0, :seller=>"0x4db96ed5e1130562", :tx_id=>"e648e442604c3e96899fa4c1b672b8add5c86c595d0199099132ae3b12e3c1e7", :tx_index=>4}, {:id=>7647917, :price=>7.0, :seller=>"0xf60cdff0ac19f0ec", :tx_id=>"2372f343732d680ec1cbe1bdf528bfd8961dc3db038d6a3bb53221d5306bd781", :tx_index=>0}, {:id=>1241712, :price=>4.0, :seller=>"0x64496f9e3a808802", :tx_id=>"903f45f72be92bf6ff58873985ac4d67a0bd8e863782c8bda300a89d3775b699", :tx_index=>1}, {:id=>12244094, :price=>5.0, :seller=>"0x85c7771750b7802a", :tx_id=>"3de7fa24bd1db9d34f619f64f7cbcd9e235870959058fddb53b3575cfb964e91", :tx_index=>0}, {:id=>2600677, :price=>4.0, :seller=>"0x923fafe9614c918f", :tx_id=>"259d070c4dae1ef55618e210e1e6cc84dbdf49ff17b98e014d79eba90deee55a", :tx_index=>2}, {:id=>5715095, :price=>3.0, :seller=>"0x2607a4ce86ed292f", :tx_id=>"9c01af83cbe58a4f51588fa01e3054dc4b1c874a09427873f5b853850900eea1", :tx_index=>11}, {:id=>10342364, :price=>7.0, :seller=>"0x57b5c6f0a4150c78", :tx_id=>"cd49469dc9e1805dbd7d4b9e8d1e5ae893265442ed8eeb4089883f78845e7149", :tx_index=>9}, {:id=>8692358, :price=>2.0, :seller=>"0xd3b951fff92d3668", :tx_id=>"dc289d5b48882a1d21f190de9a1e49565a0da68b01f80c47cc75cce9320b9d58", :tx_index=>1}, {:id=>10888458, :price=>4.0, :seller=>"0x2e15b51f44032766", :tx_id=>"d5a28c5f7dca740ccfcf320463422095cb8524780cd749a4749b660a708c32e9", :tx_index=>1}, {:id=>7432592, :price=>3.0, :seller=>"0x8485ccc82eeff2a9", :tx_id=>"f7c3aa72bc00c21a74d3dee0883960c5cbf2e5f2d0657c394a38cf70ec849a0c", :tx_index=>0}, {:id=>9935477, :price=>8.0, :seller=>"0xa732015a9089bd55", :tx_id=>"484f9e84ee0f5c0e108a86496a7a36585fd8a4497d78489c71fec23f757e0b10", :tx_index=>1}, {:id=>11454962, :price=>4.0, :seller=>"0xb537fe0409dd4936", :tx_id=>"ec964e08358f7711ea2dacb0a42484206614a48730fff52650d036a19241fad7", :tx_index=>2}, {:id=>8925456, :price=>5.0, :seller=>"0x5611b4da11d3bab6", :tx_id=>"e539cff1090a958c66d5edb8d6db2ca2b01cf6f8623ddf76daf00548deed01c8", :tx_index=>0}, {:id=>9830023, :price=>6.0, :seller=>"0x379c6680b52daba3", :tx_id=>"7d21a517f3cef62716a720a1576af06ae8184cbaba621f9e879382be070c0ec6", :tx_index=>1}, {:id=>11305425, :price=>4.0, :seller=>"0xe566f0152379ccf5", :tx_id=>"170e40f168c14832bcf841908d3847c0e6212a29e8cb227b1ffca6fad55649a3", :tx_index=>7}, {:id=>10713283, :price=>4.0, :seller=>"0xec919a6fbaa71bc9", :tx_id=>"ee9d2cacf79599b187463764d16040bfc2c2e7374a4fed2ac1859bee9a915e0d", :tx_index=>3}, {:id=>10148678, :price=>4.0, :seller=>"0x38951e1493146fe3", :tx_id=>"751894e2c6ee6b4563f8693be62fa7d846b8f708f6ecea62cf02d4080a685e3d", :tx_index=>1}]