package testassets

const ActionCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62
#,Name,,Icon,ActionCategory,,Animation{Start},VFX,Animation{End},ActionTimeline{Hit},,ClassJob,,ClassJobLevel,IsRoleAction,Range,CanTargetSelf,CanTargetParty,CanTargetFriendly,CanTargetHostile,,,TargetArea,,,,CanTargetDead,,CastType,EffectRange,XAxisModifier,,Cost{Type},Cost,,,Action{Combo},PreservesCombo,Cast<100ms>,Recast<100ms>,CooldownGroup,AttackType,Aspect,ActionProcStatus,Status{GainSelf},UnlockLink,ClassJobCategory,,,AffectsPosition,Omen,IsPvP,,,,,,,,,,,,IsPlayerAction
int32,str,bit&01,Image,ActionCategory,byte,ActionCastTimeline,ActionCastVFX,ActionTimeline,ActionTimeline,byte,ClassJob,byte,byte,bit&02,sbyte,bit&04,bit&08,bit&10,bit&20,bit&40,bit&80,bit&01,bit&02,bit&04,sbyte,bit&08,bit&10,byte,byte,byte,bit&20,byte,uint16,byte,uint16,Action,bit&40,uint16,uint16,byte,AttackType,byte,ActionProcStatus,Status,Row,ClassJobCategory,byte,bit&80,bit&01,Omen,bit&02,bit&04,bit&08,bit&10,bit&20,bit&40,bit&80,bit&01,bit&02,byte,bit&04,bit&08,bit&10
0,"",False,405,0,0,0,0,0,0,0,-1,0,0,False,0,False,False,False,False,False,False,False,False,False,0,True,True,0,0,0,False,0,0,0,0,0,False,0,0,0,0,0,0,0,0,0,0,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,False
2,"Interaction",False,405,8,0,1,0,0,1875,0,0,0,0,False,3,False,False,False,False,False,False,False,False,False,0,True,True,1,0,0,False,0,0,0,0,0,False,50,0,0,0,0,0,0,0,1,0,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,False
3,"Sprint",False,405,10,0,0,0,368,1875,0,0,0,0,False,0,True,False,False,False,False,False,False,False,False,0,True,True,1,0,0,False,0,0,0,0,0,True,0,600,56,0,0,0,0,0,1,0,False,False,0,False,False,True,False,True,True,False,False,True,0,False,True,False
4,"Mount",False,405,5,0,2,0,165,1875,0,0,0,0,False,0,True,False,False,False,False,False,False,False,False,0,True,True,1,0,0,False,0,0,0,0,0,False,10,0,0,0,0,0,0,0,1,0,False,False,0,False,False,False,False,True,True,False,False,True,0,False,True,False
5,"Teleport",False,111,10,0,1,1,164,1875,0,0,0,0,False,0,True,False,False,False,False,False,False,False,False,0,True,True,1,0,0,False,0,0,0,0,0,False,50,0,0,0,0,0,0,4,1,0,False,False,0,False,False,True,False,True,True,False,False,True,0,True,True,False
7,"attack",False,405,1,0,0,0,-1,1875,0,0,0,0,False,-1,False,False,False,True,False,False,False,False,False,0,True,True,1,0,0,True,0,0,0,0,0,True,0,0,0,-1,7,0,0,0,1,0,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,False
9,"Fast Blade",False,158,3,0,0,0,310,1875,0,1,0,1,False,-1,False,False,False,True,False,False,False,False,False,0,True,True,1,0,0,True,5,60,0,0,0,False,0,25,58,-1,7,0,0,0,38,2,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,True
11,"Savage Blade",False,157,3,0,0,0,311,1875,0,1,0,4,False,-1,False,False,False,True,False,False,False,False,False,0,True,True,1,0,0,True,5,60,0,0,9,False,0,25,58,-1,7,0,0,0,38,2,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,True
26,"Sword Oath",False,2504,2,0,29,4,442,1875,0,19,0,35,False,0,True,False,False,False,False,False,False,False,False,0,True,True,1,0,0,False,3,10,0,0,0,True,0,25,58,0,0,0,381,32,20,1,True,False,0,False,False,True,False,True,True,False,False,True,0,False,False,True
50,"Unchained",False,2554,4,0,0,0,386,1875,0,21,0,40,False,0,True,False,False,False,False,False,False,False,False,0,True,True,1,0,0,False,0,0,32,91,0,True,0,900,8,0,0,0,0,43,22,1,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,True
102,"Flaming Arrow",False,368,4,0,0,0,416,1875,0,-1,0,0,False,-1,False,False,False,False,False,False,True,False,False,0,True,True,7,5,0,False,0,0,0,0,0,True,0,600,0,-1,1,0,0,0,0,1,False,False,0,False,False,True,False,True,True,False,False,True,0,False,False,False
203,"Skyshard",False,103,9,1,52,22,2251,1883,0,0,0,0,False,25,False,False,False,False,False,False,True,False,False,0,True,True,2,8,0,True,11,0,0,0,0,False,20,0,0,8,7,0,0,0,1,0,False,False,1,False,False,True,False,True,True,True,False,True,0,False,False,False
4238,"Big Shot",False,103,9,0,54,31,3925,4012,0,0,0,0,False,30,False,False,False,True,False,False,False,False,False,0,True,True,4,30,4,True,11,0,0,0,0,False,20,0,0,8,7,0,0,0,1,0,False,False,2,False,False,True,False,True,True,True,False,True,0,False,False,False
`

const OmenCSV = `
key,0,1,2,3,4,5
#,FileName,,,,,
int32,str,str,byte,bit&01,bit&02,sbyte
0,"","",0,True,False,0
1,"general_1bf","general_1bpf",0,True,False,0
2,"general02f","general02_pf",0,True,False,0
`

const BNPCBaseCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18
#,Behavior,ActionTimelineMove,,,Scale,ModelChara,BNpcCustomize,NpcEquip,,,,ArrayEventHandler,BNpcParts,,,,,,
int32,Behavior,ActionTimelineMove,byte,byte,single,ModelChara,BNpcCustomize,NpcEquip,uint16,byte,bit&01,ArrayEventHandler,BNpcParts,bit&02,bit&04,bit&08,bit&10,byte,byte
0,0,4,0,0,1,0,0,0,0,0,False,0,0,True,True,True,True,0,0
1,0,4,0,0,1,0,2,366,0,0,False,0,0,True,True,True,True,0,0
2,0,4,0,0,1,96,0,0,0,0,False,851976,0,True,True,True,True,0,0
3,0,4,0,0,1.2,61,0,0,0,0,False,852023,0,True,True,True,True,0,0
`

const BNPCNameCSV = `
key,0,1,2,3,4,5,6,7
#,Singular,,Plural,,,,,
int32,str,sbyte,str,sbyte,sbyte,sbyte,sbyte,sbyte
0,"",0,"",0,0,1,0,0
1,"",0,"",0,0,0,0,0
2,"ruins runner",0,"ruins runners",0,0,1,0,0
3,"antelope doe",0,"antelope does",0,1,1,0,0
`

const ModelCharaCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13
#,Type,Model,Base,Variant,,,,,,,,,,
int32,byte,uint16,byte,byte,uint16,byte,bit&01,bit&02,byte,sbyte,bit&04,bit&08,bit&10,bit&20
878,3,8094,1,1,3419,1,False,False,0,0,False,False,False,False
879,3,8095,1,1,3419,1,False,False,0,0,False,False,False,False
880,3,8096,1,1,3419,1,False,False,0,0,False,False,False,False
881,3,8097,1,1,3419,1,False,False,0,0,False,False,False,False
882,3,8098,1,1,3419,1,False,False,0,0,False,False,False,False
883,3,8099,1,1,3479,1,False,False,0,0,False,False,False,False
`

const ModelSkeletonCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17
#,,,,,,,,,,,,,,,,,,
int32,single,single,single,single,uint16,uint16,uint16,uint16,uint16,uint16,uint16,uint16,single,single,uint16,byte,bit&01,byte
8094,0.2,0.5,0.3,0.5,100,300,66,300,66,300,0,0,3,3,0,0,False,1
8095,0.2,0.5,0.3,0.5,100,300,66,300,66,300,0,0,3,3,0,0,False,1
8096,0.2,0.5,0.3,0.5,100,300,66,300,66,300,0,0,3,3,0,0,False,1
8097,0.2,0.5,0.3,0.5,100,300,66,300,66,300,0,0,3,3,0,0,False,1
8098,0.2,0.5,0.3,0.5,100,300,66,300,66,300,0,0,3,3,0,0,False,1
`

const MapCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16
#,,,,Hierarchy,MapMarkerRange,Id,SizeFactor,Offset{X},Offset{Y},PlaceName{Region},PlaceName,PlaceName{Sub},DiscoveryIndex,,TerritoryType,DiscoveryArrayByte,
int32,byte,byte,sbyte,byte,uint16,str,uint16,int16,int16,PlaceName,PlaceName,PlaceName,int16,uint32,TerritoryType,bit&01,bit&02
0,0,0,0,0,0,"",100,0,0,0,0,0,-1,0,0,True,False
1,0,0,0,1,0,"default/00",100,0,0,2405,21,0,-1,0,1,True,False
2,2,2,0,1,3,"f1t1/00",200,0,0,23,52,0,-1,0,132,True,False
3,2,3,0,1,4,"f1t2/00",200,0,0,23,53,0,-1,0,133,True,False
14,3,3,1,1,60,"w1t2/01",200,0,0,24,41,373,-1,0,131,True,False
73,3,3,2,1,62,"w1t2/02",200,0,0,24,41,698,-1,0,131,True,False
178,0,0,0,1,143,"w1b4/00",200,-448,0,24,1409,0,-1,0,196,True,False
33,0,0,0,1,0,"s1fa/00",400,0,0,22,359,0,-1,0,206,True,False
403,0,0,0,1,0,"s1fa/00",400,0,0,22,359,19,-1,0,293,True,False
`

const PlaceNameCSV = `
key,0,1,2,3,4,5,6,7,8,9
#,Name,,Name{NoArticle},,,,,,,
int32,str,sbyte,str,sbyte,sbyte,sbyte,sbyte,sbyte,str,byte
0,"",0,"",0,0,1,0,0,"",0
19,"",0,"",0,0,1,0,0,"",0
21,"Eorzea",1,"Eorzea",0,0,1,0,1,"",0
41,"Ul'dah - Steps of Thal",2,"Ul'dah - Steps of Thal",0,0,1,0,1,"",0
52,"New Gridania",1,"New Gridania",0,0,1,0,1,"",0
53,"Old Gridania",1,"Old Gridania",0,0,1,0,1,"",0
359,"The Navel",1,"Navel",0,0,1,0,0,"",0
373,"Merchant Strip",2,"Merchant Strip",0,0,1,0,0,"",0
698,"Hustings Strip",2,"Hustings Strip",0,0,1,0,0,"",0
1409,"The Burning Heart",0,"Burning Heart",0,0,1,0,0,"",0
`

const TerritoryTypeCSV = `
key,0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33
#,Name,Bg,,PlaceName{Region},PlaceName{Zone},PlaceName,Map,,,TerritoryIntendedUse,,,WeatherRate,,,,,,,,,,ArrayEventHandler,,Aetheryte,,,,,,,,,
int32,str,str,byte,PlaceName,PlaceName,PlaceName,Map,byte,byte,byte,uint16,bit&01,byte,bit&02,byte,bit&04,bit&08,bit&10,bit&20,uint16,int32,int32,ArrayEventHandler,uint16,Aetheryte,int32,uint16,sbyte,bit&40,byte,byte,byte,byte,bit&80
1,"","",0,0,0,0,0,0,0,0,0,False,0,False,0,False,False,False,False,0,0,0,0,0,0,0,0,0,False,0,0,0,0,False
128,"s1t1","ffxiv/sea_s1/twn/s1t1/level/s1t1",1,22,500,28,11,2,0,0,0,False,14,True,0,True,False,False,False,1020,122007,123002,852085,0,8,-1,0,0,False,0,0,0,0,False
129,"s1t2","ffxiv/sea_s1/twn/s1t2/level/s1t2",1,22,500,29,12,2,0,0,0,False,15,True,0,True,False,False,False,1020,122007,123003,852088,0,8,-1,0,1,False,0,0,0,0,False
130,"w1t1","ffxiv/wil_w1/twn/w1t1/level/w1t1",1,24,504,40,13,4,0,0,0,False,7,True,0,True,False,False,False,1035,122008,123102,852086,0,9,-1,0,2,False,0,0,0,0,False
131,"w1t2","ffxiv/wil_w1/twn/w1t2/level/w1t2",1,24,504,41,14,4,0,0,0,False,8,True,0,True,False,False,False,1035,122008,123103,853381,0,9,-1,0,3,False,0,0,0,0,False
132,"f1t1","ffxiv/fst_f1/twn/f1t1/level/f1t1",1,23,506,52,2,3,0,0,0,False,1,True,0,True,False,False,False,1003,122009,123202,852087,0,2,-1,2,4,False,0,0,0,0,False
133,"f1t2","ffxiv/fst_f1/twn/f1t2/level/f1t2",1,23,506,53,3,3,0,0,0,False,2,True,0,True,False,False,False,1003,122009,123203,852103,0,2,-1,0,5,False,0,0,0,0,False
196,"w1b4","ffxiv/wil_w1/bah/w1b4/level/w1b4",1,24,505,1409,178,4,2,17,110,False,44,False,0,True,False,False,False,1001,122010,124534,0,0,0,-1,63,-1,False,0,0,0,0,False
206,"s1fa","ffxiv/sea_s1/fld/s1fa/level/s1fa",1,22,502,359,33,2,2,10,57,False,23,False,0,True,False,False,False,1001,-1,-1,0,0,0,-1,6,-1,False,0,0,0,0,False
293,"s1fa_2","ffxiv/sea_s1/fld/s1fa/level/s1fa",1,22,502,359,403,2,2,10,60,False,23,False,0,True,False,False,False,1001,122001,124009,0,0,0,-1,6,-1,False,0,0,0,0,False
296,"s1fa_3","ffxiv/sea_s1/fld/s1fa/level/s1fa",1,22,502,359,403,2,2,10,64,False,23,False,0,True,False,False,False,1001,122001,124010,0,0,0,-1,16,-1,False,0,0,0,0,False
`
